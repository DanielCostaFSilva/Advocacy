package auth

import (
	"context"
	"errors"
	"net/http"

	app "legalflow/internal/application/auth"
	"legalflow/internal/middleware"
)

type MeHandler struct {
	useCase GetCurrentUserUseCase
}

type GetCurrentUserUseCase interface {
	Execute(ctx context.Context, input app.GetCurrentUserInput) (*app.GetCurrentUserOutput, error)
}

func NewMeHandler(useCase GetCurrentUserUseCase) *MeHandler {
	return &MeHandler{useCase: useCase}
}

func (h *MeHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	user, err := h.useCase.Execute(r.Context(), app.GetCurrentUserInput{UserID: userID})
	if err != nil {
		if errors.Is(err, app.ErrUserNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, user)
}
