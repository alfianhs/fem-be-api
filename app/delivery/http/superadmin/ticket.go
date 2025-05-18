package superadmin_http

import (
	"app/domain/request"
	"app/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *routeSuperadmin) handleTicketRoute(route string) {
	api := h.Route.Group(route)

	api.GET("", h.Middleware.AuthSuperadmin(), h.GetTicketsList)
	api.GET("/:id", h.Middleware.AuthSuperadmin(), h.GetTicketDetail)
	api.POST("", h.Middleware.AuthSuperadmin(), h.CreateOrUpdateTicket)
	api.DELETE("/:id", h.Middleware.AuthSuperadmin(), h.DeleteTicket)
}

// GetTicketsList
//
// @Summary Get Tickets List
// @Description Get Tickets List
// @Tags Ticket-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param seriesId query string false "Series ID"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Param dir query string false "Direction asc or desc"
// @Success 200 {object} helpers.Response
// @Router /superadmin/tickets [get]
func (h *routeSuperadmin) GetTicketsList(c *gin.Context) {
	ctx := c.Request.Context()

	queryParam := c.Request.URL.Query()

	response := h.Usecase.GetTicketsList(ctx, queryParam)
	c.JSON(response.Status, response)
}

// GetTicketDetail
//
// @Summary Get Ticket Detail
// @Description Get Ticket Detail
// @Tags Ticket-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/tickets/{id} [get]
func (h *routeSuperadmin) GetTicketDetail(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	response := h.Usecase.GetTicketDetail(ctx, id)
	c.JSON(response.Status, response)
}

// CreateOrUpdateTicket
//
// @Summary Create or Update Ticket
// @Description Create or Update Ticket, if ID is empty, it will create a new ticket, otherwise it will update the existing ticket
// @Tags Ticket-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body request.TicketCreateOrUpdateRequest true "Create or Update Ticket"
// @Success 200 {object} helpers.Response
// @Router /superadmin/tickets [post]
func (h *routeSuperadmin) CreateOrUpdateTicket(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.TicketCreateOrUpdateRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	response := h.Usecase.CreateOrUpdateTicket(ctx, payload)
	c.JSON(response.Status, response)
}

// DeleteTicket
//
// @Summary Delete Ticket
// @Description Delete Ticket
// @Tags Ticket-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/tickets/{id} [delete]
func (h *routeSuperadmin) DeleteTicket(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	response := h.Usecase.DeleteTicket(ctx, id)
	c.JSON(response.Status, response)
}
