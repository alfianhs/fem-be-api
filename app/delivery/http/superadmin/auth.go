package superadmin_http

import (
	"app/domain/request"
	"app/helpers"
	jwt_helpers "app/helpers/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *routeSuperadmin) handleAuthRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.POST("/login", h.Login)
	api.GET("/profile", h.Middleware.AuthSuperadmin(), h.GetProfile)
}

// Login
//
// @Summary Login Superadmin
// @Description Login Superadmin
// @Tags Auth-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body request.SuperadminLoginRequest true "Login Superadmin"
// @Success 200 {object} helpers.Response
// @Router /superadmin/auth/login [post]
func (h *routeSuperadmin) Login(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.SuperadminLoginRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	response := h.Usecase.Login(ctx, payload)
	c.JSON(response.Status, response)
}

// GetProfile
//
// @Summary Get Profile Superadmin
// @Description Get Profile Superadmin
// @Tags Auth-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} helpers.Response
// @Router /superadmin/auth/profile [get]
func (h *routeSuperadmin) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("user_data").(jwt_helpers.SuperadminJWTClaims)

	response := h.Usecase.GetProfile(ctx, claim)
	c.JSON(response.Status, response)
}
