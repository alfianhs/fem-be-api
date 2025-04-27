package helpers

type Response struct {
	Status     int               `json:"status"`
	Message    string            `json:"message"`
	Validation map[string]string `json:"validation"`
	Data       interface{}       `json:"data"`
}

type PaginatedResponse struct {
	List  []interface{} `json:"list"`
	Limit int64         `json:"limit"`
	Page  int64         `json:"page"`
	Total int64         `json:"total"`
}

func NewResponse(status int, message string, validation map[string]string, data interface{}) Response {
	return Response{
		Status:     status,
		Message:    message,
		Validation: validation,
		Data:       data,
	}
}
