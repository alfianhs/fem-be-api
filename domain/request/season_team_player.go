package request

type SeasonTeamPlayerCreateRequest struct {
	SeasonTeamID string `form:"seasonTeamId" binding:"required"`
	PlayerID     string `form:"playerId" binding:"required"`
	Position     string `form:"position" binding:"required"`
}

type SeasonTeamPlayerUpdateRequest struct {
	Position string `form:"position" binding:"required"`
}
