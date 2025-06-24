package request

type CandidateCreateRequest struct {
	VotingID           string                      `json:"votingId"`
	SeasonTeamPlayerID string                      `json:"seasonTeamPlayerId"`
	Performance        CandidatePerformanceRequest `json:"performance"`
}

type CandidateUpdateRequest struct {
	SeasonTeamPlayerID string                      `json:"seasonTeamPlayerId"`
	Performance        CandidatePerformanceRequest `json:"performance"`
}

type CandidateVoteRequest struct {
	VotingID    string `json:"votingId"`
	CandidateID string `json:"candidateId"`
}

type CandidatePerformanceRequest struct {
	TeamLeaderboard int64 `json:"teamLeaderboard"`
	Goal            int64 `json:"goal"`
	Assist          int64 `json:"assist"`
	Save            int64 `json:"save"`
}
