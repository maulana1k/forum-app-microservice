package dto

// SignupRequest defines the body for user signup
type SignupRequest struct {
	Username string `json:"username" example:"admin"`
	Email    string `json:"email" example:"admin@example.com"`
	Password string `json:"password" example:"admin"`
}

// SigninRequest defines the body for user signin
type SigninRequest struct {
	Email    string `json:"email" example:"admin@example.com"`
	Password string `json:"password" example:"admin"`
}
