package dto

// SignupRequest defines the body for user signup
type SignupRequest struct {
	Username string `json:"username" example:"john_doe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"mypassword"`
}

// SigninRequest defines the body for user signin
type SigninRequest struct {
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"mypassword"`
}
