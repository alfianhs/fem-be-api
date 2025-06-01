package superadmin_http

import (
	"app/domain/request"
	"app/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *routeSuperadmin) handleSeasonTeamPlayerRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("", h.Middleware.AuthSuperadmin(), h.GetSeasonTeamPlayersList)
	api.GET("/position-list", h.Middleware.AuthSuperadmin(), h.GetSeasonTeamPlayersPositionList)
	api.GET("/:id", h.Middleware.AuthSuperadmin(), h.GetSeasonTeamPlayerDetail)
	api.POST("", h.Middleware.AuthSuperadmin(), h.CreateSeasonTeamPlayer)
	api.PUT("/:id", h.Middleware.AuthSuperadmin(), h.UpdateSeasonTeamPlayer)
	api.DELETE("/:id", h.Middleware.AuthSuperadmin(), h.DeleteSeasonTeamPlayer)
}

// GetSeasonTeamPlayersList
//
// @Summary Get Season Team Players List
// @Description Get Season Team Players List
// @Tags SeasonTeamPlayer-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param seasonId query string false "Season ID"
// @Param seasonTeamId query string false "Season Team ID"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Param dir query string false "Direction asc or desc"
// @Success 200 {object} helpers.Response
// @Router /superadmin/season-team-players [get]
func (h *routeSuperadmin) GetSeasonTeamPlayersList(c *gin.Context) {
	ctx := c.Request.Context()

	query := c.Request.URL.Query()

	response := h.Usecase.GetSeasonTeamPlayersList(ctx, query)
	c.JSON(response.Status, response)
}

// GetSeasonTeamPlayersPositionList
//
// @Summary Get Season Team Players Position List
// @Description Get Season Team Players Position List
// @Tags SeasonTeamPlayer-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} helpers.Response
// @Router /superadmin/season-team-players/position-list [get]
func (h *routeSuperadmin) GetSeasonTeamPlayersPositionList(c *gin.Context) {
	ctx := c.Request.Context()

	response := h.Usecase.GetPlayerPositionsList(ctx)
	c.JSON(response.Status, response)
}

// GetSeasonTeamPlayerDetail
//
// @Summary Get Season Team Player Detail
// @Description Get Season Team Player Detail
// @Tags SeasonTeamPlayer-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Season Team Player ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/season-team-players/{id} [get]
func (h *routeSuperadmin) GetSeasonTeamPlayerDetail(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	response := h.Usecase.GetSeasonTeamPlayerDetail(ctx, id)
	c.JSON(response.Status, response)
}

// CreateSeasonTeamPlayer
//
// @Summary Create Season Team Player
// @Description Create Season Team Player
// @Tags SeasonTeamPlayer-Superadmin
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param payload formData request.SeasonTeamPlayerCreateRequest true "Create Season Team Player"
// @Param image formData file true "Image"
// @Success 201 {object} helpers.Response
// @Router /superadmin/season-team-players [post]
func (h *routeSuperadmin) CreateSeasonTeamPlayer(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.SeasonTeamPlayerCreateRequest{}
	err := c.ShouldBind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	response := h.Usecase.CreateSeasonTeamPlayer(ctx, payload, c.Request)
	c.JSON(response.Status, response)
}

// UpdateSeasonTeamPlayer
//
// @Summary Update Season Team Player
// @Description Update Season Team Player
// @Tags SeasonTeamPlayer-Superadmin
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Season Team Player ID"
// @Param payload formData request.SeasonTeamPlayerUpdateRequest true "Update Season Team Player"
// @Param image formData file false "Image"
// @Success 200 {object} helpers.Response
// @Router /superadmin/season-team-players/{id} [put]
func (h *routeSuperadmin) UpdateSeasonTeamPlayer(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.SeasonTeamPlayerUpdateRequest{}
	err := c.ShouldBind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	id := c.Param("id")

	response := h.Usecase.UpdateSeasonTeamPlayer(ctx, id, payload, c.Request)
	c.JSON(response.Status, response)
}

// DeleteSeasonTeamPlayer
//
// @Summary Delete Season Team Player
// @Description Delete Season Team Player
// @Tags SeasonTeamPlayer-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Season Team Player ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/season-team-players/{id} [delete]
func (h *routeSuperadmin) DeleteSeasonTeamPlayer(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	response := h.Usecase.DeleteSeasonTeamPlayer(ctx, id)
	c.JSON(response.Status, response)
}
