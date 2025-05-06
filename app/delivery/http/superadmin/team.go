package superadmin_http

import (
	"app/domain/request"
	"app/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *routeSuperadmin) handleTeamRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("", h.Middleware.AuthSuperadmin(), h.GetTeamsList)
	api.GET("/:id", h.Middleware.AuthSuperadmin(), h.GetTeamDetail)
	api.POST("", h.Middleware.AuthSuperadmin(), h.CreateTeam)
	api.PUT("/:id", h.Middleware.AuthSuperadmin(), h.UpdateTeam)
	api.DELETE("/:id", h.Middleware.AuthSuperadmin(), h.DeleteTeam)
}

// GetTeamsList
//
// @Summary Get Teams List
// @Description Get Teams List
// @Tags Team-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} helpers.Response
// @Router /superadmin/teams [get]
func (h *routeSuperadmin) GetTeamsList(c *gin.Context) {
	ctx := c.Request.Context()

	query := c.Request.URL.Query()

	response := h.Usecase.GetTeamsList(ctx, query)
	c.JSON(response.Status, response)
}

// GetTeamDetail
//
// @Summary Get Team Detail
// @Description Get Team Detail
// @Tags Team-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/teams/{id} [get]
func (h *routeSuperadmin) GetTeamDetail(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	response := h.Usecase.GetTeamDetail(ctx, id)
	c.JSON(response.Status, response)
}

// CreateTeam
//
// @Summary Create Team
// @Description Create Team
// @Tags Team-Superadmin
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param payload formData request.TeamCreateRequest true "Create Team"
// @Param logo formData file true "Logo"
// @Success 201 {object} helpers.Response
// @Router /superadmin/teams [post]
func (h *routeSuperadmin) CreateTeam(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.TeamCreateRequest{}
	err := c.ShouldBind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	response := h.Usecase.CreateTeam(ctx, payload, c.Request)
	c.JSON(response.Status, response)
}

// UpdateTeam
//
// @Summary Update Team
// @Description Update Team
// @Tags Team-Superadmin
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Team ID"
// @Param payload formData request.TeamUpdateRequest true "Update Team"
// @Param logo formData file false "Logo"
// @Success 200 {object} helpers.Response
// @Router /superadmin/teams/{id} [put]
func (h *routeSuperadmin) UpdateTeam(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.TeamUpdateRequest{}
	err := c.ShouldBind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	id := c.Param("id")
	request := c.Request

	response := h.Usecase.UpdateTeam(ctx, id, payload, request)
	c.JSON(response.Status, response)
}

// DeleteTeam
//
// @Summary Delete Team
// @Description Delete Team
// @Tags Team-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/teams/{id} [delete]
func (h *routeSuperadmin) DeleteTeam(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	response := h.Usecase.DeleteTeam(ctx, id)
	c.JSON(response.Status, response)
}
