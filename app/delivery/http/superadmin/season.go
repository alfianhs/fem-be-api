package superadmin_http

import (
	"app/domain/request"
	"app/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *routeSuperadmin) handleSeasonRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("", h.Middleware.AuthSuperadmin(), h.GetSeasonsList)
	api.GET("/:id", h.Middleware.AuthSuperadmin(), h.GetSeasonDetail)
	api.POST("", h.Middleware.AuthSuperadmin(), h.CreateSeason)
	api.PUT("/:id", h.Middleware.AuthSuperadmin(), h.UpdateSeason)
	api.DELETE("/:id", h.Middleware.AuthSuperadmin(), h.DeleteSeason)
	api.PUT("/:id/status", h.Middleware.AuthSuperadmin(), h.UpdateSeasonStatus)
}

// GetSeasonsList
//
// @Summary Get Seasons
// @Description Get Seasons
// @Tags Season-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Param dir query string false "Direction asc or desc"
// @Success 200 {object} helpers.Response
// @Router /superadmin/seasons [get]
func (h *routeSuperadmin) GetSeasonsList(c *gin.Context) {
	ctx := c.Request.Context()

	query := c.Request.URL.Query()

	response := h.Usecase.GetSeasonsList(ctx, query)
	c.JSON(response.Status, response)
}

// GetSeasonDetail
//
// @Summary Get Season Detail
// @Description Get Season Detail
// @Tags Season-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Season ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/seasons/{id} [get]
func (h *routeSuperadmin) GetSeasonDetail(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	response := h.Usecase.GetSeasonDetail(ctx, id)
	c.JSON(response.Status, response)
}

// CreateSeason
//
// @Summary Create Season
// @Description Create Season
// @Tags Season-Superadmin
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param payload formData request.SeasonCreateRequest true "Create Season"
// @Param logo formData file true "Logo"
// @Param banner formData file true "Banner"
// @Success 201 {object} helpers.Response
// @Router /superadmin/seasons [post]
func (h *routeSuperadmin) CreateSeason(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.SeasonCreateRequest{}
	err := c.ShouldBind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	response := h.Usecase.CreateSeason(ctx, payload, c.Request)
	c.JSON(response.Status, response)
}

// UpdateSeason
//
// @Summary Update Season
// @Description Update Season
// @Tags Season-Superadmin
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Season ID"
// @Param payload formData request.SeasonUpdateRequest true "Update Season"
// @Param logo formData file false "Logo"
// @Param banner formData file false "Banner"
// @Success 200 {object} helpers.Response
// @Router /superadmin/seasons/{id} [put]
func (h *routeSuperadmin) UpdateSeason(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.SeasonUpdateRequest{}
	err := c.ShouldBind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	id := c.Param("id")

	response := h.Usecase.UpdateSeason(ctx, id, payload, c.Request)
	c.JSON(response.Status, response)
}

// DeleteSeason
//
// @Summary Delete Season
// @Description Delete Season
// @Tags Season-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Season ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/seasons/{id} [delete]
func (h *routeSuperadmin) DeleteSeason(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	response := h.Usecase.DeleteSeason(ctx, id)
	c.JSON(response.Status, response)
}

// UpdateSeasonStatus
//
// @Summary Update Season Status
// @Description Update Season Status
// @Tags Season-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Season ID"
// @Param payload body request.SeasonStatusUpdateRequest true "Update Season Status"
// @Success 200 {object} helpers.Response
// @Router /superadmin/seasons/{id}/status [put]
func (h *routeSuperadmin) UpdateSeasonStatus(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.SeasonStatusUpdateRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	id := c.Param("id")

	response := h.Usecase.UpdateSeasonStatus(ctx, id, payload)
	c.JSON(response.Status, response)
}
