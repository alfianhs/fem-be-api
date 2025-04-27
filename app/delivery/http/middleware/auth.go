package middleware

import (
	"app/helpers"
	jwt_helpers "app/helpers/jwt"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (m *appMiddleware) AuthSuperadmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token from header
		requestToken := c.Request.Header.Get("Authorization")
		if requestToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(
				http.StatusUnauthorized,
				"Unauthorized: Missing Authorization Header",
				nil,
				nil,
			))
			return
		}

		// prepend prefix "Bearer "
		if !strings.HasPrefix(requestToken, "Bearer ") {
			requestToken = "Bearer " + requestToken
		}

		// check token format
		splitToken := strings.Split(requestToken, "Bearer ")
		if len(splitToken) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(
				http.StatusUnauthorized,
				"Unauthorized: Invalid Token Format",
				nil,
				nil,
			))
			return
		}

		// get token from split
		tokenString := splitToken[1]

		// validate token
		token, err := jwt.ParseWithClaims(tokenString, &jwt_helpers.SuperadminJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secretKeySuperadmin), nil
		})

		// check if token is valid
		if !token.Valid {
			if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(
					http.StatusUnauthorized,
					"Unauthorized: Invalid Token Signature",
					nil,
					nil,
				))
				return
			}

			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(
					http.StatusUnauthorized,
					"Unauthorized: Token Expired",
					nil,
					nil,
				))
				return
			}
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(
				http.StatusUnauthorized,
				err.Error(),
				nil,
				nil,
			))
			return
		}

		claims, ok := token.Claims.(*jwt_helpers.SuperadminJWTClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(
				http.StatusUnauthorized,
				"Unauthorized: Invalid Token Claims",
				nil,
				nil,
			))
			return
		}

		// set claims to context
		c.Set("user_data", *claims)
		c.Next()
	}
}

func (m *appMiddleware) AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token from header
		requestToken := c.Request.Header.Get("Authorization")
		if requestToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(
				http.StatusUnauthorized,
				"Unauthorized: Missing Authorization Header",
				nil,
				nil,
			))
			return
		}

		// prepend prefix "Bearer "
		if !strings.HasPrefix(requestToken, "Bearer ") {
			requestToken = "Bearer " + requestToken
		}

		// check token format
		splitToken := strings.Split(requestToken, "Bearer ")
		if len(splitToken) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(
				http.StatusUnauthorized,
				"Unauthorized: Invalid Token Format",
				nil,
				nil,
			))
			return
		}

		// get token from split
		tokenString := splitToken[1]

		// validate token
		token, err := jwt.ParseWithClaims(tokenString, &jwt_helpers.AdminJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secretKeyAdmin), nil
		})

		// check if token is valid
		if !token.Valid {
			if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(
					http.StatusUnauthorized,
					"Unauthorized: Invalid Token Signature",
					nil,
					nil,
				))
				return
			}

			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(
					http.StatusUnauthorized,
					"Unauthorized: Token Expired",
					nil,
					nil,
				))
				return
			}
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(
				http.StatusUnauthorized,
				err.Error(),
				nil,
				nil,
			))
			return
		}

		claims, ok := token.Claims.(*jwt_helpers.AdminJWTClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(
				http.StatusUnauthorized,
				"Unauthorized: Invalid Token Claims",
				nil,
				nil,
			))
			return
		}

		// set claims to context
		c.Set("user_data", *claims)
		c.Next()
	}
}

func (m *appMiddleware) AuthMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token from header
		requestToken := c.Request.Header.Get("Authorization")
		if requestToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(
				http.StatusUnauthorized,
				"Unauthorized: Missing Authorization Header",
				nil,
				nil,
			))
			return
		}

		// prepend prefix "Bearer "
		if !strings.HasPrefix(requestToken, "Bearer ") {
			requestToken = "Bearer " + requestToken
		}

		// check token format
		splitToken := strings.Split(requestToken, "Bearer ")
		if len(splitToken) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(
				http.StatusUnauthorized,
				"Unauthorized: Invalid Token Format",
				nil,
				nil,
			))
			return
		}

		// get token from split
		tokenString := splitToken[1]

		// validate token
		token, err := jwt.ParseWithClaims(tokenString, &jwt_helpers.MemberJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secretKeyMember), nil
		})

		// check if token is valid
		if !token.Valid {
			if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(
					http.StatusUnauthorized,
					"Unauthorized: Invalid Token Signature",
					nil,
					nil,
				))
				return
			}

			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(
					http.StatusUnauthorized,
					"Unauthorized: Token Expired",
					nil,
					nil,
				))
				return
			}
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(
				http.StatusUnauthorized,
				err.Error(),
				nil,
				nil,
			))
			return
		}

		claims, ok := token.Claims.(*jwt_helpers.MemberJWTClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(
				http.StatusUnauthorized,
				"Unauthorized: Invalid Token Claims",
				nil,
				nil,
			))
			return
		}

		// set claims to context
		c.Set("user_data", *claims)
		c.Next()
	}
}
