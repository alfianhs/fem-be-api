package member_http

import (
	"app/domain/request"
	"app/helpers"
	jwt_helpers "app/helpers/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *routeMember) handleAuthRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.POST("/register", h.Register)
	api.POST("/verify-email", h.VerifyEmail)
	api.POST("/resend-email-verification", h.ResendEmailVerification)
	api.POST("/login", h.Login)
	api.GET("/profile", h.Middleware.AuthMember(), h.GetProfile)
}

// Register
//
// @Summary Register Member
// @Description Register Member
// @Tags Auth-Member
// @Accept json
// @Produce json
// @Param payload body request.MemberRegisterRequest true "Register Member"
// @Success 201 {object} helpers.Response
// @Router /member/auth/register [post]
func (h *routeMember) Register(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.MemberRegisterRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	response := h.Usecase.Register(ctx, payload)
	c.JSON(response.Status, response)
}

// VerifyEmail
//
// @Summary Verify Email Member
// @Description Verify Email Member
// @Tags Auth-Member
// @Accept json
// @Produce json
// @Param payload body request.VerifyEmailRequest true "Verify Email Member"
// @Success 200 {object} helpers.Response
// @Router /member/auth/verify-email [post]
func (h *routeMember) VerifyEmail(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.VerifyEmailRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	response := h.Usecase.VerifyEmail(ctx, payload)
	c.JSON(response.Status, response)
}

// ResendVerifyEmail
//
// @Summary Resend Verify Email Member
// @Description Resend Verify Email Member
// @Tags Auth-Member
// @Accept json
// @Produce json
// @Param payload body request.ResendEmailVerificationRequest true "Resend Verify Email Member"
// @Success 200 {object} helpers.Response
// @Router /member/auth/resend-email-verification [post]
func (h *routeMember) ResendEmailVerification(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.ResendEmailVerificationRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	response := h.Usecase.ResendEmailVerification(ctx, payload)
	c.JSON(response.Status, response)
}

// Login
//
// @Summary Login Member
// @Description Login Member
// @Tags Auth-Member
// @Accept json
// @Produce json
// @Param payload body request.MemberLoginRequest true "Login Member"
// @Success 200 {object} helpers.Response
// @Router /member/auth/login [post]
func (h *routeMember) Login(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.MemberLoginRequest{}
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
// @Summary Get Profile Member
// @Description Get Profile Member
// @Tags Auth-Member
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} helpers.Response
// @Router /member/auth/profile [get]
func (h *routeMember) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("user_data").(jwt_helpers.MemberJWTClaims)

	response := h.Usecase.GetProfile(ctx, claim)
	c.JSON(response.Status, response)
}
