package request

type VenueCreateRequest struct {
	Name string `json:"name"`
}

type VenueUpdateRequest struct {
	Name string `json:"name"`
}
