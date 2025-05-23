package superadmin_usecase

import (
	"app/domain/request"
	"app/helpers"
	jwt_helpers "app/helpers/jwt"
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (u *superadminAppUsecase) Login(ctx context.Context, payload request.SuperadminLoginRequest) helpers.Response {
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

	// check superadmin
	superadmin, err := u.mongoDbRepo.FetchOneSuperadmin(ctx, map[string]interface{}{
		"email": payload.Email,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if superadmin == nil {
		return helpers.NewResponse(http.StatusBadRequest, "User not found", nil, nil)
	}

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(superadmin.Password), []byte(payload.Password)); err != nil {
		return helpers.NewResponse(http.StatusBadRequest, "Wrong password", nil, nil)
	}

	// generate token
	now := time.Now()
	logrus.Info(time.Duration(jwt_helpers.GetJWTTTL()))
	token, err := jwt_helpers.GenerateJWTTokenSuperadmin(jwt_helpers.SuperadminJWTClaims{
		UserID: superadmin.ID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    "superadmin",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwt_helpers.GetJWTTTL()) * time.Minute)),
		},
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Login successful", nil, map[string]any{
		"token": token,
		"user":  superadmin,
	})
}

func (u *superadminAppUsecase) GetProfile(ctx context.Context, claim jwt_helpers.SuperadminJWTClaims) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check superadmin
	superadmin, err := u.mongoDbRepo.FetchOneSuperadmin(ctx, map[string]interface{}{
		"id": claim.UserID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if superadmin == nil {
		return helpers.NewResponse(http.StatusBadRequest, "User not found", nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Get profile success", nil, superadmin)
}
