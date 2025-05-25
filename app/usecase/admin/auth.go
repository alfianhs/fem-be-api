package admin_usecase

import (
	"app/domain/request"
	"app/helpers"
	jwt_helpers "app/helpers/jwt"
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (u *adminAppUsecase) Login(ctx context.Context, payload request.AdminLoginRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	errValidation := make(map[string]string)
	if payload.Username == "" {
		errValidation["username"] = "Username field is required"
	}
	if payload.Password == "" {
		errValidation["password"] = "Password field is required"
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	// check admin
	admin, err := u.mongoDbRepo.FetchOneAdmin(ctx, map[string]interface{}{
		"username": payload.Username,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if admin == nil {
		return helpers.NewResponse(http.StatusBadRequest, "User not found", nil, nil)
	}

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(payload.Password)); err != nil {
		return helpers.NewResponse(http.StatusBadRequest, "Wrong password", nil, nil)
	}

	// generate token
	now := time.Now()
	expiredAt := now.Add(time.Duration(jwt_helpers.GetJWTTTL()) * time.Minute)
	token, err := jwt_helpers.GenerateJWTTokenAdmin(jwt_helpers.AdminJWTClaims{
		UserID: admin.ID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    "admin",
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
		"user":      admin,
	})
}

func (u *adminAppUsecase) GetProfile(ctx context.Context, claim jwt_helpers.AdminJWTClaims) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check admin
	admin, err := u.mongoDbRepo.FetchOneAdmin(ctx, map[string]interface{}{
		"id": claim.UserID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if admin == nil {
		return helpers.NewResponse(http.StatusBadRequest, "User not found", nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Get profile success", nil, admin)
}
