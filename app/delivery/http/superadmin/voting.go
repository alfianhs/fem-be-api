package superadmin_http

import (
	"app/domain/request"
	"app/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *routeSuperadmin) handleVotingRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("", h.Middleware.AuthSuperadmin(), h.GetVotingList)
	api.GET("/:id", h.Middleware.AuthSuperadmin(), h.GetVotingDetail)
	api.POST("", h.Middleware.AuthSuperadmin(), h.CreateVoting)
	api.PUT("/:id", h.Middleware.AuthSuperadmin(), h.UpdateVoting)
	api.DELETE("/:id", h.Middleware.AuthSuperadmin(), h.DeleteVoting)
}

// GetVotingList
//
// @Summary Get Voting List
// @Description Get list of all votings
// @Tags Voting-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort field"
// @Param dir query string false "Direction asc or desc"
// @Param status query int false "Status filter"
// @Success 200 {object} helpers.Response
// @Router /superadmin/votings [get]
func (h *routeSuperadmin) GetVotingList(c *gin.Context) {
	ctx := c.Request.Context()
	queryParam := c.Request.URL.Query()

	response := h.Usecase.GetVotingList(ctx, queryParam)
	c.JSON(response.Status, response)
}

// GetVotingDetail
//
// @Summary Get Voting Detail
// @Description Get detail of a single voting by ID
// @Tags Voting-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Voting ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/votings/{id} [get]
func (h *routeSuperadmin) GetVotingDetail(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	response := h.Usecase.GetVotingDetail(ctx, id)
	c.JSON(response.Status, response)
}

// CreateVoting
//
// @Summary Create Voting
// @Description Create a new voting
// @Tags Voting-Superadmin
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param payload formData request.VotingCreateRequest true "Create Voting"
// @Param banner formData file true "Banner"
// @Success 201 {object} helpers.Response
// @Router /superadmin/votings [post]
func (h *routeSuperadmin) CreateVoting(c *gin.Context) {
	ctx := c.Request.Context()
	var payload request.VotingCreateRequest

	if err := c.ShouldBind(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(
			http.StatusBadRequest, "Invalid request payload", nil, nil,
		))
		return
	}

	response := h.Usecase.CreateVoting(ctx, payload, c.Request)
	c.JSON(response.Status, response)
}

// UpdateVoting
//
// @Summary Update Voting
// @Description Update a voting by ID
// @Tags Voting-Superadmin
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Voting ID"
// @Param payload formData request.VotingUpdateRequest false "Update Voting"
// @Param banner formData file false "Banner"
// @Success 200 {object} helpers.Response
// @Router /superadmin/votings/{id} [put]
func (h *routeSuperadmin) UpdateVoting(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request.VotingUpdateRequest{}
	err := c.ShouldBind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	id := c.Param("id")

	response := h.Usecase.UpdateVoting(ctx, id, payload, c.Request)
	c.JSON(response.Status, response)
}

// DeleteVoting
//
// @Summary Delete Voting
// @Description Soft-delete a voting by ID
// @Tags Voting-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Voting ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/votings/{id} [delete]
func (h *routeSuperadmin) DeleteVoting(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	response := h.Usecase.DeleteVoting(ctx, id)
	c.JSON(response.Status, response)
}
