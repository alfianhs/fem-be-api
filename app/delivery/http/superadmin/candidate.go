package superadmin_http

import (
	"app/domain/request"
	"app/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *routeSuperadmin) handleCandidateRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("", h.Middleware.AuthSuperadmin(), h.GetCandidateList)
	api.GET("/:id", h.Middleware.AuthSuperadmin(), h.GetCandidateDetail)
	api.POST("", h.Middleware.AuthSuperadmin(), h.CreateCandidate)
	api.PUT("/:id", h.Middleware.AuthSuperadmin(), h.UpdateCandidate)
	api.DELETE("/:id", h.Middleware.AuthSuperadmin(), h.DeleteCandidate)
}

// GetCandidateList
//
// @Summary Get Candidate List
// @Description Get list of all candidates
// @Tags Candidate-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort field"
// @Param dir query string false "Direction asc or desc"
// @Param votingId query string false "Filter by Voting ID"
// @Param seasonTeamPlayerId query int false "Filter by SeasonTeamPlayer ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/candidates [get]
func (h *routeSuperadmin) GetCandidateList(c *gin.Context) {
	ctx := c.Request.Context()
	queryParam := c.Request.URL.Query()

	response := h.Usecase.GetCandidateList(ctx, queryParam)
	c.JSON(response.Status, response)
}

// GetCandidateDetail
//
// @Summary Get Candidate Detail
// @Description Get detail of a single candidate by ID
// @Tags Candidate-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Candidate ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/candidates/{id} [get]
func (h *routeSuperadmin) GetCandidateDetail(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	response := h.Usecase.GetCandidateDetail(ctx, id)
	c.JSON(response.Status, response)
}

// CreateCandidate
//
// @Summary Create Candidate
// @Description Create a new candidate
// @Tags Candidate-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body request.CandidateCreateRequest true "Create Candidate"
// @Success 201 {object} helpers.Response
// @Router /superadmin/candidates [post]
func (h *routeSuperadmin) CreateCandidate(c *gin.Context) {
	ctx := c.Request.Context()

	var payload request.CandidateCreateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(
			http.StatusBadRequest, "Invalid request payload", nil, nil,
		))
		return
	}

	response := h.Usecase.CreateCandidate(ctx, payload)
	c.JSON(response.Status, response)
}

// UpdateCandidate
//
// @Summary Update Candidate
// @Description Update an existing candidate by ID
// @Tags Candidate-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Candidate ID"
// @Param payload body request.CandidateUpdateRequest false "Update Candidate"
// @Success 200 {object} helpers.Response
// @Router /superadmin/candidates/{id} [put]
func (h *routeSuperadmin) UpdateCandidate(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	var payload request.CandidateUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(
			http.StatusBadRequest, "Invalid request payload", nil, nil,
		))
		return
	}

	response := h.Usecase.UpdateCandidate(ctx, id, payload)
	c.JSON(response.Status, response)
}

// DeleteCandidate
//
// @Summary Delete Candidate
// @Description Soft-delete a candidate by ID
// @Tags Candidate-Superadmin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Candidate ID"
// @Success 200 {object} helpers.Response
// @Router /superadmin/candidates/{id} [delete]
func (h *routeSuperadmin) DeleteCandidate(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	response := h.Usecase.DeleteCandidate(ctx, id)
	c.JSON(response.Status, response)
}
