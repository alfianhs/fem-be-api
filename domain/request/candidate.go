package request

type CandidateCreateRequest struct {
	VotingID           string `json:"votingId" binding:"required"`
	SeasonTeamPlayerID string `json:"seasonTeamPlayerId" binding:"required"`
	Performance        string `json:"performance" binding:"required"`
}

type CandidateUpdateRequest struct {
	SeasonTeamPlayerID string `json:"seasonTeamPlayerId"`
	Performance        string `json:"performance"`
}

type CandidateVoteRequest struct {
	VotingID    string `json:"votingId" binding:"required"`
	CandidateID string `json:"candidateId" binding:"required"`
}
