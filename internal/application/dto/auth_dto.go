package dto

type AuthResult struct {
	Token    string `json:"token"`
	UserID   string `json:"userId"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}