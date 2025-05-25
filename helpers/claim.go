package helpers

import (
	jwt_helpers "app/helpers/jwt"

	"github.com/gin-gonic/gin"
)

func GetClaim(c *gin.Context) jwt_helpers.MemberJWTClaims {
	val, exists := c.Get("user_data")
	if !exists {
		return jwt_helpers.MemberJWTClaims{}
	}

	claim, ok := val.(jwt_helpers.MemberJWTClaims)
	if !ok {
		return jwt_helpers.MemberJWTClaims{}
	}

	return claim
}
