package utils

// Meta struktur metadata untuk pagination
type Meta struct {
	CurrentPage  int `json:"current_page"`
	LimitPerPage int `json:"limit_per_page"`
	TotalData    int `json:"total_data"`
	TotalPages   int `json:"total_pages"`
}

// JSendResponse struktur standar balikan API sesuai pedoman HMIK
type JSendResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    *Meta       `json:"meta,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
}

// BuildSuccessResponse menghasilkan response sukses biasa
func BuildSuccessResponse(message string, data interface{}) JSendResponse {
	return JSendResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

// BuildFailResponse menghasilkan response gagal (client error)
func BuildFailResponse(message string, errors []string) JSendResponse {
	return JSendResponse{
		Status:  "fail",
		Message: message,
		Data:    nil,
		Errors:  errors,
	}
}

// BuildErrorResponse menghasilkan response error (server error)
func BuildErrorResponse(message string) JSendResponse {
	return JSendResponse{
		Status:  "error",
		Message: message,
		Data:    nil,
	}
}
