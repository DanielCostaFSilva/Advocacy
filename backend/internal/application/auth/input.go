package auth

type AuthenticateUserInput struct {
	Email    string
	Password string
}

type GetCurrentUserInput struct {
	UserID string
}
