package request

type TeamCreateRequest struct {
	Name string `form:"name" validate:"required"`
}

type TeamUpdateRequest struct {
	Name string `form:"name" validate:"required"`
}
