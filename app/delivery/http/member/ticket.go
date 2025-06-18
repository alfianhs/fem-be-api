package member_http

import (
	"github.com/gin-gonic/gin"
)

func (h *routeMember) handleTicketRoute(route string) {
	api := h.Route.Group(route)

	api.GET("", h.GetTicketsList)
	api.GET("/:id", h.GetTicketDetail)
}

// GetMemberTicketsList
// @Summary Get Tickets List
// @Description Get Tickets List
// @Tags Ticket-Member
// @Accept json
// @Produce json
// @Param seriesId query string false "Series ID"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Param dir query string false "Direction asc or desc"
// @Success 200 {object} helpers.Response
// @Router /member/tickets [get]
func (h *routeMember) GetTicketsList(c *gin.Context) {
	ctx := c.Request.Context()

	queryParam := c.Request.URL.Query()

	response := h.Usecase.GetTicketsList(ctx, queryParam)
	c.JSON(response.Status, response)
}

// GetMemberTicketDetail
// @Summary Get Ticket Detail
// @Description Get Ticket Detail
// @Tags Ticket-Member
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Success 200 {object} helpers.Response
// @Router /member/tickets/{id} [get]
func (h *routeMember) GetTicketDetail(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	response := h.Usecase.GetTicketDetail(ctx, id)
	c.JSON(response.Status, response)
}
