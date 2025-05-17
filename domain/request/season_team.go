package request

type SeasonTeamCreateRequest struct {
	TeamIds []string `json:"teamIds"`
}

type SeasonTeamManageRequest struct {
	AddedTeamIds   []string `json:"addedTeamIds"`
	RemovedTeamIds []string `json:"removedTeamIds"`
}
