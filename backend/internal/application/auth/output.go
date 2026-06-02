package auth

type AuthenticateUserOutput struct {
	UserID string
	Name   string
	Email  string
}

type GetCurrentUserOutput struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
