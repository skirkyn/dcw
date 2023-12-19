package dto

type ErrorResponse struct {
	Err string `json:"err"`
}

func NewErrorResponse(str string) Response[string] {
	return ErrorResponse{str}
}
func (e ErrorResponse) Error() string {
	return e.Err
}
