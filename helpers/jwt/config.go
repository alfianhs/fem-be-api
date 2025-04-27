package jwt_helpers

import (
	"os"
	"strconv"
)

func GetJWTSecretKeySuperadmin() string {
	return os.Getenv("JWT_SECRET_KEY_SUPERADMIN")
}

func GetJWTSecretKeyAdmin() string {
	return os.Getenv("JWT_SECRET_KEY_ADMIN")
}

func GetJWTSecretKeyMember() string {
	return os.Getenv("JWT_SECRET_KEY_MEMBER")
}

func GetJWTTTL() int {
	ttl, _ := strconv.Atoi(os.Getenv("JWT_TTL"))
	return ttl
}
