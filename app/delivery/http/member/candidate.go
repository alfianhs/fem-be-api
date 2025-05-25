package member_http

import (
	"app/domain/request"
	"app/helpers"
	jwt_helpers "app/helpers/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *routeMember) handleCandidateRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("", h.Middleware.OptionalAuthMember(), h.GetCandidateList)
	api.POST("/vote", h.Middleware.AuthMember(), h.CandidateVote)
}

// GetCandidateList
// @Summary Get Candidate List
// @Description Get list of all candidates
// @Tags Candidate-Member
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort field"
// @Param dir query string false "Direction asc or desc"
// @Param votingId query string false "Filter by Voting ID"
// @Param seasonTeamPlayerId query int false "Filter by SeasonTeamPlayer ID"
// @Success 200 {object} helpers.Response
// @Router /member/candidates [get]
func (h *routeMember) GetCandidateList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := helpers.GetClaim(c)

	response := h.Usecase.GetCandidateList(ctx, claim, c.Request.URL.Query())
	c.JSON(response.Status, response)
}

// CandidateVote
//
// @Summary Candidate Vote
// @Description Candidate Vote
// @Tags Candidate-Member
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body request.CandidateVoteRequest true "Candidate Vote"
// @Success 200 {object} helpers.Response
// @Router /member/candidates/vote [post]
func (h *routeMember) CandidateVote(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("user_data").(jwt_helpers.MemberJWTClaims)
	var payload request.CandidateVoteRequest
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid json data", nil, nil))
		return
	}

	response := h.Usecase.CandidateVote(ctx, claim, payload)
	c.JSON(response.Status, response)
}
