package request

type SeasonTeamPlayerCreateRequest struct {
	SeasonTeamID string `form:"seasonTeamId" validate:"required"`
	PlayerID     string `form:"playerId" validate:"required"`
	Position     string `form:"position" validate:"required"`
}

type SeasonTeamPlayerUpdateRequest struct {
	Position string `form:"position" validate:"required"`
}
