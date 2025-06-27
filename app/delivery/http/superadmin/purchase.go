package superadmin_http

import "github.com/gin-gonic/gin"

func (h *routeSuperadmin) handlePurchaseRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("", h.Middleware.AuthSuperadmin(), h.GetPurchasesList)
}

// GetPurchasesList
//
//	@Summary		Get Purchases List
//	@Description	Get Purchases List
//	@Tags			Purchase-Superadmin
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			search	query	string	false	"Search by invoice external id / email"
//	@Param			page	query	int	false	"Page"
//	@Param			limit	query	int	false	"Limit"
//	@Param			sort	query	string	false	"Sort"
//	@Param			dir		query	string	false	"Direction asc or desc"
//	@Success		200		{object}	helpers.Response
//	@Router			/superadmin/purchases [get]
func (h *routeSuperadmin) GetPurchasesList(c *gin.Context) {
	ctx := c.Request.Context()

	query := c.Request.URL.Query()

	response := h.Usecase.GetPurchasesList(ctx, query)
	c.JSON(response.Status, response)
}
