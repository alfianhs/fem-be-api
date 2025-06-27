package superadmin_http

import "github.com/gin-gonic/gin"

func (h *routeSuperadmin) handleDashboardRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("", h.Middleware.AuthSuperadmin(), h.GetDashboard)
}

// GetDashboard
//
//	@Summary		Get Dashboard
//	@Description	Get Dashboard
//	@Tags			Dashboard-Superadmin
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			seasonId	query	string	false	"Season ID"
//	@Success		200		{object}	helpers.Response
//	@Router			/superadmin/dashboard [get]
func (h *routeSuperadmin) GetDashboard(c *gin.Context) {
	ctx := c.Request.Context()

	query := c.Request.URL.Query()

	response := h.Usecase.GetDashboard(ctx, query)
	c.JSON(response.Status, response)
}
