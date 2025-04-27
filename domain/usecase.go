package domain

import (
	"app/domain/request"
	"app/helpers"
	jwt_helpers "app/helpers/jwt"
	"context"
)

type SuperadminAppUsecase interface {
	// Auth
	Login(ctx context.Context, payload request.SuperadminLoginRequest) helpers.Response
	GetProfile(ctx context.Context, claim jwt_helpers.SuperadminJWTClaims) helpers.Response
}

type AdminAppUsecase interface {
	// Auth
	Login(ctx context.Context, payload request.AdminLoginRequest) helpers.Response
	GetProfile(ctx context.Context, claim jwt_helpers.AdminJWTClaims) helpers.Response
}

type MemberAppUsecase interface {
	// Auth
	Register(ctx context.Context, payload request.MemberRegisterRequest) helpers.Response
	VerifyEmail(ctx context.Context, payload request.VerifyEmailRequest) helpers.Response
	ResendEmailVerification(ctx context.Context, payload request.ResendEmailVerificationRequest) helpers.Response
	Login(ctx context.Context, payload request.MemberLoginRequest) helpers.Response
	GetProfile(ctx context.Context, claim jwt_helpers.MemberJWTClaims) helpers.Response
}
