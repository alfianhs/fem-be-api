package member_http

import (
	jwt_helpers "app/helpers/jwt"

	"github.com/gin-gonic/gin"
)

func (h *routeMember) handleTicketPurchaseRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("", h.Middleware.AuthMember(), h.GetTicketPurchasesList)
}

// GetTicketPurchasesList
//
//	@Summary		Get ticket purchases list
//	@Description	Get ticket purchases list
//	@Tags			TicketPurchase-Member
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			page	query	int	false	"Page"
//	@Param			limit	query	int	false	"Limit"
//	@Param			sort	query	string	false	"Sort"
//	@Param			dir		query	string	false	"Direction asc or desc"
//	@Success		200		{object}	helpers.Response
//	@Router			/member/ticket-purchases [get]
func (h *routeMember) GetTicketPurchasesList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("user_data").(jwt_helpers.MemberJWTClaims)
	queryParam := c.Request.URL.Query()

	response := h.Usecase.GetTicketPurchasesList(ctx, claim, queryParam)
	c.JSON(response.Status, response)
}
