package member_http

import (
	"github.com/gin-gonic/gin"
)

func (h *routeMember) handleSeriesRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("", h.GetSeriesList)
	api.GET("/:id", h.GetSeriesDetail)
}

// GetSeriesList
//
//	@Summary Get Series
//	@Description Get Series
//	@Tags Series-Member
//	@Accept json
//	@Produce json
//	@Param search query string false "Search by name"
//	@Param page query int false "Page"
//	@Param limit query int false "Limit"
//	@Param sort query string false "Sort"
//	@Param dir query string false "Direction asc or desc"
//	@Success 200 {object} helpers.Response
//	@Router /member/series [get]
func (h *routeMember) GetSeriesList(c *gin.Context) {
	ctx := c.Request.Context()

	queryParam := c.Request.URL.Query()

	response := h.Usecase.GetSeriesList(ctx, queryParam)
	c.JSON(response.Status, response)
}

// GetSeriesDetail
//
//	@Summary Get Series Detail
//	@Description Get Series Detail
//	@Tags Series-Member
//	@Accept json
//	@Produce json
//	@Param id path string true "Series ID"
//	@Success 200 {object} helpers.Response
//	@Router /member/series/{id} [get]
func (h *routeMember) GetSeriesDetail(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	response := h.Usecase.GetSeriesDetail(ctx, id)
	c.JSON(response.Status, response)
}
