package pkgutil

type HTTPResponse struct {
	Code    int    `json:"-"`
	Message string `json:"message,omitempty" example:"Success"`
	Data    any    `json:"data,omitempty" `
	Meta    any    `json:"meta,omitempty" `
}

type PaginationResponse struct {
	TotalData int `json:"total_data" example:"1"`
	TotalPage int `json:"total_page" example:"1"`
	Page      int `json:"page" example:"1"`
	Limit     int `json:"limit" example:"10"`
	Data      any `json:"data,omitempty" `
}

type ErrValidationResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type MetaResponse struct {
	Total  int `json:"total" example:"1"`
	Offset int `json:"offset" example:"0"`
	Limit  int `json:"limit" example:"10"`
}
