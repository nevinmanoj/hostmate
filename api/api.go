package api

type GetAllResponsePage[T any] struct {
	Message      string `json:"message"`
	StatusCode   int    `json:"status_code"`
	TotalRecords int64  `json:"total_records"`
	PageSize     int    `json:"page_size"`
	CurrentPage  int    `json:"current_page"`
	TotalPages   int    `json:"total_pages"`
	Data         []T    `json:"data"`
}
type GetResponsePage[T any] struct {
	Message    string `json:"message"`
	Data       T      `json:"data"`
	StatusCode int    `json:"status_code"`
}

type PostResponsePage[T any] struct {
	Message    string `json:"message"`
	Data       T      `json:"data"`
	StatusCode int    `json:"status_code"`
}
type PutResponsePage[T any] struct {
	Message    string `json:"message"`
	Data       T      `json:"data"`
	StatusCode int    `json:"status_code"`
}
type DeleteResponsePage struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}
type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
	User  any    `json:"user"`
}
