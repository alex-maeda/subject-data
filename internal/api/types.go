// Package api defines the request/response types for the Subject Property API.
package api

// ErrorResponse is the body of a 4xx/5xx response.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail is the error payload. Field is optional and points at the
// offending request path when relevant (e.g. "filters[2].type").
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}
