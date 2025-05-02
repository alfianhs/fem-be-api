package superadmin_http

import (
	"app/domain/request"
	"app/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *routeSuperadmin) handleVenueRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("", h.Middleware.AuthSuperadmin(), h.GetVenuesList)
	api.GET("/:id", h.Middleware.AuthSuperadmin(), h.GetVenueDetail)
	api.POST("", h.Middleware.AuthSuperadmin(), h.CreateVenue)
	api.PUT("/:id", h.Middleware.AuthSuperadmin(), h.UpdateVenue)
	api.DELETE("/:id", h.Middleware.AuthSuperadmin(), h.DeleteVenue)
}

// GetVenuesList
//
// @Summary Get Venues List
// @Description Get Venues List
// @Tags Venue-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Param dir query string false "Direction asc or desc"
// @Success 200 {object} helpers.Response
// @Router /superadmin/venues [get]
func (h *routeSuperadmin) GetVenuesList(c *gin.Context) {
	ctx := c.Request.Context()

	options := map[string]interface{}{
		"query": c.Request.URL.Query(),
	}

	response := h.Usecase.GetVenueList(ctx, options)
	c.JSON(response.Status, response)
}

// GetVenueDetail
//
// @Summary Get Venue Detail
// @Description Get Venue Detail
// @Tags Venue-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Venue ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/venues/{id} [get]
func (h *routeSuperadmin) GetVenueDetail(c *gin.Context) {
	ctx := c.Request.Context()

	options := map[string]interface{}{
		"id": c.Param("id"),
	}

	response := h.Usecase.GetVenueDetail(ctx, options)
	c.JSON(response.Status, response)
}

// CreateVenue
//
// @Summary Create Venue
// @Description Create Venue
// @Tags Venue-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body request.VenueCreateRequest true "Create Venue"
// @Success 201 {object} helpers.Response
// @Router /superadmin/venues [post]
func (h *routeSuperadmin) CreateVenue(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.VenueCreateRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	response := h.Usecase.CreateVenue(ctx, payload)
	c.JSON(response.Status, response)
}

// UpdateVenue
//
// @Summary Update Venue
// @Description Update Venue
// @Tags Venue-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Venue ID"
// @Param payload body request.VenueUpdateRequest true "Update Venue"
// @Success 200 {object} helpers.Response
// @Router /superadmin/venues/{id} [put]
func (h *routeSuperadmin) UpdateVenue(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.VenueUpdateRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	options := map[string]interface{}{
		"id":      c.Param("id"),
		"payload": payload,
	}

	response := h.Usecase.UpdateVenue(ctx, options)
	c.JSON(response.Status, response)
}

// DeleteVenue
//
// @Summary Delete Venue
// @Description Delete Venue
// @Tags Venue-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Venue ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/venues/{id} [delete]
func (h *routeSuperadmin) DeleteVenue(c *gin.Context) {
	ctx := c.Request.Context()

	options := map[string]interface{}{
		"id": c.Param("id"),
	}

	response := h.Usecase.DeleteVenue(ctx, options)
	c.JSON(response.Status, response)
}
