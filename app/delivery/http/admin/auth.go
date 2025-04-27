package admin_http

import (
	"app/domain/request"
	"app/helpers"
	jwt_helpers "app/helpers/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *routeAdmin) handleAuthRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.POST("/login", h.Login)
	api.GET("/profile", h.Middleware.AuthAdmin(), h.GetProfile)
}

// Login
//
// @Summary Login Admin
// @Description Login Admin
// @Tags Auth-Admin
// @Accept json
// @Produce json
// @Param payload body request.AdminLoginRequest true "Login Admin"
// @Success 200 {object} helpers.Response
// @Router /admin/auth/login [post]
func (h *routeAdmin) Login(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.AdminLoginRequest{}
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
// @Summary Get Profile Admin
// @Description Get Profile Admin
// @Tags Auth-Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} helpers.Response
// @Router /admin/auth/profile [get]
func (h *routeAdmin) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("user_data").(jwt_helpers.AdminJWTClaims)

	response := h.Usecase.GetProfile(ctx, claim)
	c.JSON(response.Status, response)
}
