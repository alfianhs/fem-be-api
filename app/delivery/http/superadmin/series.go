package superadmin_http

import (
	"app/domain/request"
	"app/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *routeSuperadmin) handleSeriesRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("", h.Middleware.AuthSuperadmin(), h.GetSeriesList)
	api.GET("/:id", h.Middleware.AuthSuperadmin(), h.GetSeriesDetail)
	api.POST("", h.Middleware.AuthSuperadmin(), h.CreateSeries)
	api.PUT("/:id", h.Middleware.AuthSuperadmin(), h.UpdateSeries)
	api.DELETE("/:id", h.Middleware.AuthSuperadmin(), h.DeleteSeries)
}

// GetSeriesList
//
// @Summary Get Series
// @Description Get Series
// @Tags Series-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search by name"
// @Param seasonId query string false "Season ID"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Param dir query string false "Direction asc or desc"
// @Success 200 {object} helpers.Response
// @Router /superadmin/series [get]
func (h *routeSuperadmin) GetSeriesList(c *gin.Context) {
	ctx := c.Request.Context()

	queryParam := c.Request.URL.Query()

	response := h.Usecase.GetSeriesList(ctx, queryParam)
	c.JSON(response.Status, response)
}

// GetSeriesDetail
//
// @Summary Get Series Detail
// @Description Get Series Detail
// @Tags Series-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Series ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/series/{id} [get]
func (h *routeSuperadmin) GetSeriesDetail(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	response := h.Usecase.GetSeriesDetail(ctx, id)
	c.JSON(response.Status, response)
}

// CreateSeries
//
// @Summary Create Series
// @Description Create Series
// @Tags Series-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body request.SeriesCreateRequest true "Create Series"
// @Success 201 {object} helpers.Response
// @Router /superadmin/series [post]
func (h *routeSuperadmin) CreateSeries(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.SeriesCreateRequest{}
	err := c.ShouldBind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	response := h.Usecase.CreateSeries(ctx, payload)
	c.JSON(response.Status, response)
}

// UpdateSeries
//
// @Summary Update Series
// @Description Update Series
// @Tags Series-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Series ID"
// @Param payload body request.SeriesUpdateRequest true "Update Series"
// @Success 200 {object} helpers.Response
// @Router /superadmin/series/{id} [put]
func (h *routeSuperadmin) UpdateSeries(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.SeriesUpdateRequest{}
	err := c.ShouldBind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	id := c.Param("id")

	response := h.Usecase.UpdateSeries(ctx, id, payload)
	c.JSON(response.Status, response)
}

// DeleteSeries
//
// @Summary Delete Series
// @Description Delete Series
// @Tags Series-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Series ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/series/{id} [delete]
func (h *routeSuperadmin) DeleteSeries(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	response := h.Usecase.DeleteSeries(ctx, id)
	c.JSON(response.Status, response)
}
