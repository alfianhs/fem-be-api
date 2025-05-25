package middleware

import (
	jwt_helpers "app/helpers/jwt"
	"io"

	"github.com/gin-gonic/gin"
)

type appMiddleware struct {
	secretKeySuperadmin string
	secretKeyAdmin      string
	secretKeyMember     string
}

func NewAppMiddleware() AppMiddleware {
	return &appMiddleware{
		secretKeySuperadmin: jwt_helpers.GetJWTSecretKeySuperadmin(),
		secretKeyAdmin:      jwt_helpers.GetJWTSecretKeyAdmin(),
		secretKeyMember:     jwt_helpers.GetJWTSecretKeyMember(),
	}
}

type AppMiddleware interface {
	AuthSuperadmin() gin.HandlerFunc
	AuthAdmin() gin.HandlerFunc
	AuthMember() gin.HandlerFunc
	OptionalAuthMember() gin.HandlerFunc
	Logger(writer io.Writer) gin.HandlerFunc
	Recovery() gin.HandlerFunc
}
