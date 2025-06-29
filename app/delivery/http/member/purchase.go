package member_http

import (
	"app/domain/request"
	"app/helpers"
	jwt_helpers "app/helpers/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *routeMember) handlePurchaseRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("/:id", h.Middleware.AuthMember(), h.GetPurchaseDetail)
	api.POST("", h.Middleware.AuthMember(), h.CreatePurchase)
	api.POST("/packages", h.Middleware.AuthMember(), h.CreatePackagePurchase)
}

// GetPurchaseDetail
//
//	@Summary		Get purchase detail
//	@Description	Get purchase detail
//	@Tags			Purchase-Member
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id path string true "Purchase ID"
//	@Success		200		{object}	helpers.Response
//	@Router			/member/purchases/{id} [get]
func (h *routeMember) GetPurchaseDetail(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")
	claim := c.MustGet("user_data").(jwt_helpers.MemberJWTClaims)

	response := h.Usecase.GetPurchaseDetail(ctx, claim, id)
	c.JSON(response.Status, response)
}

// CreatePurchase
//
//	@Summary		Create purchase
//	@Description	Create purchase
//	@Tags			Purchase-Member
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body	request.CreatePurchaseRequest	true	"Create purchase"
//	@Success		200		{object}	helpers.Response
//	@Router			/member/purchases [post]
func (h *routeMember) CreatePurchase(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("user_data").(jwt_helpers.MemberJWTClaims)
	var payload request.CreatePurchaseRequest
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	response := h.Usecase.CreatePurchase(ctx, claim, payload)
	c.JSON(response.Status, response)
}

// CreatePackagePurchase
//
//	@Summary		Create package purchase
//	@Description	Create package purchase
//	@Tags			Purchase-Member
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body	request.CreatePurchaseRequest	true	"Create package purchase"
//	@Success		200		{object}	helpers.Response
//	@Router			/member/purchases/packages [post]
func (h *routeMember) CreatePackagePurchase(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("user_data").(jwt_helpers.MemberJWTClaims)
	var payload request.CreatePurchaseRequest
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	response := h.Usecase.CreatePackagePurchase(ctx, claim, payload)
	c.JSON(response.Status, response)
}
