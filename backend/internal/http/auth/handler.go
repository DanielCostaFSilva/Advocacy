package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	app "legalflow/internal/application/auth"
)

type AuthenticateUserUseCase interface {
	Execute(ctx context.Context, input app.AuthenticateUserInput) (*app.AuthenticateUserOutput, error)
}

type TokenProvider interface {
	Generate(userID, email string) (string, error)
}

type LoginHandler struct {
	authUseCase AuthenticateUserUseCase
	tokenProv   TokenProvider
}

func NewLoginHandler(authUseCase AuthenticateUserUseCase, tokenProv TokenProvider) *LoginHandler {
	return &LoginHandler{authUseCase: authUseCase, tokenProv: tokenProv}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, LoginErrorResponse{Error: "invalid request body"})
		return
	}

	if req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, LoginErrorResponse{Error: "email and password are required"})
		return
	}

	user, err := h.authUseCase.Execute(r.Context(), app.AuthenticateUserInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, app.ErrInvalidCredentials) {
			writeJSON(w, http.StatusUnauthorized, LoginErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, app.ErrInvalidInput) {
			writeJSON(w, http.StatusBadRequest, LoginErrorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, LoginErrorResponse{Error: "internal server error"})
		return
	}

	token, err := h.tokenProv.Generate(user.UserID, user.Email)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, LoginErrorResponse{Error: "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		User: UserResponse{
			ID:    user.UserID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
