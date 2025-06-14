package member_usecase

import (
	mongo_model "app/domain/model/mongo"
	"app/domain/request"
	"app/helpers"
	jwt_helpers "app/helpers/jwt"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (u *memberAppUsecase) Register(ctx context.Context, payload request.MemberRegisterRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	errValidation := make(map[string]string)
	if payload.Name == "" {
		errValidation["name"] = "Name field is required"
	}
	if payload.Email == "" {
		errValidation["email"] = "Email field is required"
	} else {
		if !helpers.IsValidEmail(payload.Email) {
			errValidation["email"] = "Invalid email format"
		}
	}
	if payload.Password == "" {
		errValidation["password"] = "Password field is required"
	} else {
		if !helpers.IsValidLengthPassword(payload.Password) {
			errValidation["password"] = "Password must be at least 8 characters"
		}
		if !helpers.IsStrongPassword(payload.Password) {
			errValidation["password"] = "Password must contain at least one uppercase letter, one lowercase letter, and one number"
		}
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	// check member
	member, err := u.mongoDbRepo.FetchOneMember(ctx, map[string]interface{}{
		"email": payload.Email,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if member != nil {
		return helpers.NewResponse(http.StatusBadRequest, "User already exist", nil, nil)
	}

	// hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	// generate email token
	emailToken, _ := helpers.GenerateSecureRandomChar(64)

	now := time.Now()
	// create member
	member = &mongo_model.Member{
		ID:         primitive.NewObjectID(),
		Name:       payload.Name,
		Email:      payload.Email,
		Password:   string(hashedPassword),
		EmailToken: emailToken,
		IsVerified: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// save
	if err := u.mongoDbRepo.CreateOneMember(ctx, member); err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	// send email verification
	go sendEmailVerification(member)

	return helpers.NewResponse(
		http.StatusCreated,
		"Register Successful",
		nil,
		member,
	)
}

func sendEmailVerification(member *mongo_model.Member) {
	// helper email
	mailer := helpers.NewSMTPMailer()

	// set verification link
	baseFeUrl := helpers.GetFEUrl()
	verificationLink := fmt.Sprintf("%s/verify-email/%s?email=%s", baseFeUrl, member.EmailToken, member.Email)

	// get template
	subject, body := helpers.GetEmailVerificationTemplate()

	// replace string template
	dataReplace := map[string]string{
		"link_verification": verificationLink,
		"user_name":         member.Name,
	}
	finalBody := helpers.StringReplacer(body, dataReplace)

	// setup mail content
	mailer.To([]string{member.Email})
	mailer.Subject(subject)
	mailer.Body(finalBody)

	// send
	if err := mailer.Send(); err != nil {
		logrus.Errorf("Send Email to %s error %v", member.Email, err)
	}
}

func (u *memberAppUsecase) VerifyEmail(ctx context.Context, payload request.VerifyEmailRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	errValidation := make(map[string]string)
	if payload.Token == "" {
		errValidation["token"] = "Token field is required"
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	// check member
	member, err := u.mongoDbRepo.FetchOneMember(ctx, map[string]interface{}{
		"emailToken": payload.Token,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if member == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Token is invalid", nil, nil)
	}
	if member.IsVerified {
		return helpers.NewResponse(http.StatusBadRequest, "Email already verified", nil, nil)
	}

	// update member
	now := time.Now()
	member.IsVerified = true
	member.VerifiedAt = &now
	if err := u.mongoDbRepo.UpdatePartialMember(ctx, map[string]interface{}{
		"emailToken": payload.Token,
	}, map[string]interface{}{
		"isVerified": member.IsVerified,
		"verifiedAt": member.VerifiedAt,
		"emailToken": "",
	}); err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Email verification successful", nil, member)
}

func (u *memberAppUsecase) ResendEmailVerification(ctx context.Context, payload request.ResendEmailVerificationRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	errValidation := make(map[string]string)
	if payload.Email == "" {
		errValidation["email"] = "Email field is required"
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	// check member
	member, err := u.mongoDbRepo.FetchOneMember(ctx, map[string]interface{}{
		"email": payload.Email,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if member == nil {
		return helpers.NewResponse(http.StatusBadRequest, "User not found", nil, nil)
	}
	if member.IsVerified {
		return helpers.NewResponse(http.StatusBadRequest, "Email already verified", nil, nil)
	}

	// update email token member
	emailToken, _ := helpers.GenerateSecureRandomChar(64)
	member.EmailToken = emailToken

	// save
	if err := u.mongoDbRepo.UpdatePartialMember(ctx, map[string]interface{}{
		"id": member.ID,
	}, map[string]interface{}{
		"emailToken": member.EmailToken,
	}); err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	// send email verification
	go sendEmailVerification(member)

	return helpers.NewResponse(http.StatusOK, "Resend email verification successful", nil, nil)
}

func (u *memberAppUsecase) Login(ctx context.Context, payload request.MemberLoginRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	errValidation := make(map[string]string)
	if payload.Email == "" {
		errValidation["email"] = "Email field is required"
	}
	if payload.Password == "" {
		errValidation["password"] = "Password field is required"
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	// check member
	member, err := u.mongoDbRepo.FetchOneMember(ctx, map[string]interface{}{
		"email": payload.Email,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if member == nil {
		return helpers.NewResponse(http.StatusBadRequest, "User not found", nil, nil)
	}
	if !member.IsVerified {
		return helpers.NewResponse(http.StatusBadRequest, "Email not verified", nil, nil)
	}

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(payload.Password)); err != nil {
		return helpers.NewResponse(http.StatusBadRequest, "Wrong password", nil, nil)
	}

	// generate token
	now := time.Now()
	expiredAt := now.Add(time.Duration(jwt_helpers.GetJWTTTL()) * time.Minute)
	token, err := jwt_helpers.GenerateJWTTokenMember(jwt_helpers.MemberJWTClaims{
		UserID: member.ID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    "member",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiredAt),
		},
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Login successful", nil, map[string]any{
		"token":     token,
		"expiredAt": expiredAt,
		"user":      member,
	})
}

func (u *memberAppUsecase) GetProfile(ctx context.Context, claim jwt_helpers.MemberJWTClaims) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check member
	member, err := u.mongoDbRepo.FetchOneMember(ctx, map[string]interface{}{
		"id": claim.UserID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if member == nil {
		return helpers.NewResponse(http.StatusBadRequest, "User not found", nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Get profile success", nil, member)
}
