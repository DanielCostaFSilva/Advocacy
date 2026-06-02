package legalcase

import (
	"time"

	"github.com/google/uuid"
)

type CaseStatus string

const (
	CaseStatusDraft     CaseStatus = "draft"
	CaseStatusActive    CaseStatus = "active"
	CaseStatusSuspended CaseStatus = "suspended"
	CaseStatusClosed    CaseStatus = "closed"
)

type Case struct {
	ID          uuid.UUID
	ClientID    uuid.UUID
	Number      string
	Title       string
	Description string
	Court       string
	Status      CaseStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewCase(clientID uuid.UUID, number, title, description, court string) (*Case, error) {
	if clientID == uuid.Nil {
		return nil, ErrInvalidClientID
	}

	if len(number) < 5 {
		return nil, ErrInvalidCaseNumber
	}

	if len(title) < 5 {
		return nil, ErrInvalidCaseTitle
	}

	if court == "" {
		return nil, ErrInvalidCourt
	}

	now := time.Now()

	return &Case{
		ID:          uuid.New(),
		ClientID:    clientID,
		Number:      number,
		Title:       title,
		Description: description,
		Court:       court,
		Status:      CaseStatusDraft,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
