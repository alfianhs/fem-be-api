package member_http

import (
	"github.com/gin-gonic/gin"
)

func (h *routeMember) handleVotingRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("", h.GetVotingList)
}

// GetVotingList
//
// @Summary Get Voting List
// @Description Get list of all votings
// @Tags Voting-Member
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort field"
// @Param dir query string false "Direction asc or desc"
// @Param seriesId query string false "Filter by Series ID"
// @Param status query int false "Status filter"
// @Success 200 {object} helpers.Response
// @Router /member/votings [get]
func (h *routeMember) GetVotingList(c *gin.Context) {
	ctx := c.Request.Context()
	queryParam := c.Request.URL.Query()

	response := h.Usecase.GetVotingList(ctx, queryParam)
	c.JSON(response.Status, response)
}
