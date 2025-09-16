package dto

type ErrorResponse struct {
	Error   string `json:"message"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}
