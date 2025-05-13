package superadmin_http

import (
	"app/domain/request"
	"app/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *routeSuperadmin) handleSeasonTeamRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("", h.Middleware.AuthSuperadmin(), h.GetSeasonTeamsList)
	api.GET("/:id", h.Middleware.AuthSuperadmin(), h.GetSeasonTeamDetail)
	api.POST("", h.Middleware.AuthSuperadmin(), h.CreateSeasonTeam)
	api.DELETE("/:id", h.Middleware.AuthSuperadmin(), h.DeleteSeasonTeam)
}

// GetSeasonTeamsList
//
// @Summary Get Season Teams List
// @Description Get Season Teams List
// @Tags SeasonTeam-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param seasonId query string false "Season ID"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Param dir query string false "Direction asc or desc"
// @Success 200 {object} helpers.Response
// @Router /superadmin/season-teams [get]
func (h *routeSuperadmin) GetSeasonTeamsList(c *gin.Context) {
	ctx := c.Request.Context()

	query := c.Request.URL.Query()

	response := h.Usecase.GetSeasonTeamsList(ctx, query)
	c.JSON(response.Status, response)
}

// GetSeasonTeamDetail
//
// @Summary Get Season Team Detail
// @Description Get Season Team Detail
// @Tags SeasonTeam-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Season Team ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/season-teams/{id} [get]
func (h *routeSuperadmin) GetSeasonTeamDetail(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	response := h.Usecase.GetSeasonTeamDetail(ctx, id)
	c.JSON(response.Status, response)
}

// CreateSeasonTeam
//
// @Summary Create Season Team
// @Description Create Season Team
// @Tags SeasonTeam-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body request.SeasonTeamCreateRequest true "Create Season Team"
// @Success 201 {object} helpers.Response
// @Router /superadmin/season-teams [post]
func (h *routeSuperadmin) CreateSeasonTeam(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.SeasonTeamCreateRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	response := h.Usecase.CreateSeasonTeam(ctx, payload)
	c.JSON(response.Status, response)
}

// DeleteSeasonTeam
//
// @Summary Delete Season Team
// @Description Delete Season Team
// @Tags SeasonTeam-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Season Team ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/season-teams/{id} [delete]
func (h *routeSuperadmin) DeleteSeasonTeam(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	response := h.Usecase.DeleteSeasonTeam(ctx, id)
	c.JSON(response.Status, response)
}
