package jwt_helpers

import "github.com/golang-jwt/jwt/v5"

type SuperadminJWTClaims struct {
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}

type AdminJWTClaims struct {
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}

type MemberJWTClaims struct {
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}
