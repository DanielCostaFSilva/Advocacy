package hearing

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	app "legalflow/internal/application/hearing"
)

type ListHearingsUseCase interface {
	Execute(ctx context.Context, input app.ListHearingsInput) (*app.ListHearingsOutput, error)
}

type ListHandler struct {
	useCase ListHearingsUseCase
}

func NewListHandler(useCase ListHearingsUseCase) *ListHandler {
	return &ListHandler{useCase: useCase}
}

func (h *ListHandler) Register(r chi.Router) {
	r.Get("/hearings", h.List)
}

func (h *ListHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	input := app.ListHearingsInput{
		Page:      page,
		PageSize:  pageSize,
		CaseID:    r.URL.Query().Get("case_id"),
		Type:      r.URL.Query().Get("type"),
		StartDate: r.URL.Query().Get("start_date"),
		EndDate:   r.URL.Query().Get("end_date"),
		Sort:      r.URL.Query().Get("sort"),
		Order:     r.URL.Query().Get("order"),
	}

	output, err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	resp := ListHearingsResponse{
		Data: make([]HearingItem, len(output.Data)),
		Pagination: PaginationInfo{
			Page:       output.Page,
			PageSize:   output.PageSize,
			Total:      output.Total,
			TotalPages: output.TotalPages,
		},
	}
	for i, h := range output.Data {
		resp.Data[i] = HearingItem{
			ID:          h.ID,
			CaseID:      h.CaseID,
			Title:       h.Title,
			Type:        h.Type,
			Location:    h.Location,
			ScheduledAt: h.ScheduledAt,
			CreatedAt:   h.CreatedAt,
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

type ListHearingsResponse struct {
	Data       []HearingItem  `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

type HearingItem struct {
	ID          string    `json:"id"`
	CaseID      string    `json:"case_id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Location    string    `json:"location"`
	ScheduledAt time.Time `json:"scheduled_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type PaginationInfo struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}
