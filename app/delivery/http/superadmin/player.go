package superadmin_http

import (
	"app/domain/request"
	"app/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *routeSuperadmin) handlePlayerRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("", h.Middleware.AuthSuperadmin(), h.GetPlayersList)
	api.GET("/:id", h.Middleware.AuthSuperadmin(), h.GetPlayerDetail)
	api.POST("", h.Middleware.AuthSuperadmin(), h.CreatePlayer)
	api.PUT("/:id", h.Middleware.AuthSuperadmin(), h.UpdatePlayer)
	api.DELETE("/:id", h.Middleware.AuthSuperadmin(), h.DeletePlayer)
}

// GetPlayersList
//
// @Summary Get Players List
// @Description Get Players List
// @Tags Player-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Param dir query string false "Direction asc or desc"
// @Success 200 {object} helpers.Response
// @Router /superadmin/players [get]
func (h *routeSuperadmin) GetPlayersList(c *gin.Context) {
	ctx := c.Request.Context()

	query := c.Request.URL.Query()

	response := h.Usecase.GetPlayerList(ctx, query)
	c.JSON(response.Status, response)
}

// GetPlayerDetail
//
// @Summary Get Player Detail
// @Description Get Player Detail
// @Tags Player-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Player ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/players/{id} [get]
func (h *routeSuperadmin) GetPlayerDetail(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	response := h.Usecase.GetPlayerDetail(ctx, id)
	c.JSON(response.Status, response)
}

// CreatePlayer
//
// @Summary Create Player
// @Description Create Player
// @Tags Player-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body request.PlayerCreateRequest true "Create Player"
// @Success 201 {object} helpers.Response
// @Router /superadmin/players [post]
func (h *routeSuperadmin) CreatePlayer(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.PlayerCreateRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	response := h.Usecase.CreatePlayer(ctx, payload)
	c.JSON(response.Status, response)
}

// UpdatePlayer
//
// @Summary Update Player
// @Description Update Player
// @Tags Player-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Player ID"
// @Param payload body request.PlayerUpdateRequest true "Update Player"
// @Success 200 {object} helpers.Response
// @Router /superadmin/players/{id} [put]
func (h *routeSuperadmin) UpdatePlayer(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.PlayerUpdateRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	id := c.Param("id")

	response := h.Usecase.UpdatePlayer(ctx, id, payload)
	c.JSON(response.Status, response)
}

// DeletePlayer
//
// @Summary Delete Player
// @Description Delete Player
// @Tags Player-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Player ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/players/{id} [delete]
func (h *routeSuperadmin) DeletePlayer(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	response := h.Usecase.DeletePlayer(ctx, id)
	c.JSON(response.Status, response)
}
