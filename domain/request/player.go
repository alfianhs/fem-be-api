package request

type PlayerCreateRequest struct {
	Name      string  `json:"name"`
	StageName *string `json:"stageName"`
}

type PlayerUpdateRequest struct {
	Name      string  `json:"name"`
	StageName *string `json:"stageName"`
}
