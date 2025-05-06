package request

type PlayerCreateRequest struct {
	Name string `json:"name"`
}

type PlayerUpdateRequest struct {
	Name string `json:"name"`
}
