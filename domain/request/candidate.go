package request

type CandidateCreateRequest struct {
	VotingID           string `json:"votingId"`
	SeasonTeamPlayerID string `json:"seasonTeamPlayerId"`
	Performance        string `json:"performance"`
}

type CandidateUpdateRequest struct {
	SeasonTeamPlayerID string `json:"seasonTeamPlayerId"`
	Performance        string `json:"performance"`
}

type CandidateVoteRequest struct {
	VotingID    string `json:"votingId"`
	CandidateID string `json:"candidateId"`
}
