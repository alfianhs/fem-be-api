package request

type TeamCreateRequest struct {
	Name string `form:"name" binding:"required"`
}

type TeamUpdateRequest struct {
	Name string `form:"name" binding:"required"`
}
