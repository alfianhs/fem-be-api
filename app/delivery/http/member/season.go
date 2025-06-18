package member_http

import (
	"github.com/gin-gonic/gin"
)

func (h *routeMember) handleSeasonRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("/active", h.GetActiveSeasonDetail)
}

// GetMemberActiveSeasonDetail
// @Summary Get Active Season Detail
// @Description Get Active Season Detail
// @Tags Season-Member
// @Accept json
// @Produce json
// @Success 200 {object} helpers.Response
// @Router /member/seasons/active [get]
func (h *routeMember) GetActiveSeasonDetail(c *gin.Context) {
	ctx := c.Request.Context()

	response := h.Usecase.GetActiveSeasonDetail(ctx)
	c.JSON(response.Status, response)
}
