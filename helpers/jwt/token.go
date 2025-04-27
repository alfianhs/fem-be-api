package jwt_helpers

import (
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTTokenSuperadmin(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(GetJWTSecretKeySuperadmin()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateJWTTokenAdmin(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(GetJWTSecretKeyAdmin()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateJWTTokenMember(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(GetJWTSecretKeyMember()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
