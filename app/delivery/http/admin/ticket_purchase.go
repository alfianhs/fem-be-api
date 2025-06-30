package admin_http

import (
	"app/domain/request"
	"app/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *routeAdmin) handleTicketPurchaseRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("/used-today", h.Middleware.AuthAdmin(), h.GetTicketPurchasesListIsUsedToday)
	api.POST("/scan", h.Middleware.AuthAdmin(), h.Scan)
}

// GetTicketPurchasesListIsUsedToday
//
// @Summary Get Ticket Purchases List Is Used Today
// @Description Get Ticket Purchases List Is Used Today
// @Tags TicketPurchase-Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} helpers.Response
// @Router /admin/ticket-purchases/used-today [get]
func (h *routeAdmin) GetTicketPurchasesListIsUsedToday(c *gin.Context) {
	ctx := c.Request.Context()

	response := h.Usecase.GetListTicketPurchasesIsUsedToday(ctx)
	c.JSON(response.Status, response)
}

// Scan
//
// @Summary Scan Ticket Purchase
// @Description Scan Ticket Purchase
// @Tags TicketPurchase-Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body request.ScanTicketPurchaseRequest true "Scan Ticket Purchase"
// @Success 200 {object} helpers.Response
// @Router /admin/ticket-purchases/scan [post]
func (h *routeAdmin) Scan(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.ScanTicketPurchaseRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	response := h.Usecase.ScanTicketPurchase(ctx, payload)
	c.JSON(response.Status, response)
}
