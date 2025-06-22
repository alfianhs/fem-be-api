package middleware

import (
	jwt_helpers "app/helpers/jwt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

type appMiddleware struct {
	secretKeySuperadmin string
	secretKeyAdmin      string
	secretKeyMember     string
	xenditCallbackToken string
}

func NewAppMiddleware() AppMiddleware {
	return &appMiddleware{
		secretKeySuperadmin: jwt_helpers.GetJWTSecretKeySuperadmin(),
		secretKeyAdmin:      jwt_helpers.GetJWTSecretKeyAdmin(),
		secretKeyMember:     jwt_helpers.GetJWTSecretKeyMember(),
		xenditCallbackToken: os.Getenv("XENDIT_CALLBACK_TOKEN"),
	}
}

type AppMiddleware interface {
	AuthSuperadmin() gin.HandlerFunc
	AuthAdmin() gin.HandlerFunc
	AuthMember() gin.HandlerFunc
	OptionalAuthMember() gin.HandlerFunc
	AuthXendit() gin.HandlerFunc
	Logger(writer io.Writer) gin.HandlerFunc
	Recovery() gin.HandlerFunc
}
