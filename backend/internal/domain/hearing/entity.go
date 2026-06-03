package hearing

import (
	"time"

	"github.com/google/uuid"
)

type Hearing struct {
	ID          uuid.UUID
	CaseID      uuid.UUID
	Title       string
	Description string
	Type        HearingType
	Location    string
	ScheduledAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewHearing(caseID uuid.UUID, title, description string, hearingType HearingType, location string, scheduledAt time.Time) (*Hearing, error) {
	if caseID == uuid.Nil {
		return nil, ErrInvalidCaseID
	}

	if len(title) < 5 {
		return nil, ErrInvalidTitle
	}

	switch hearingType {
	case HearingTypeConciliation, HearingTypeInstruction, HearingTypeJudgment, HearingTypeVirtual:
	default:
		return nil, ErrInvalidHearingType
	}

	if location == "" {
		return nil, ErrInvalidLocation
	}

	if scheduledAt.Before(time.Now()) {
		return nil, ErrInvalidScheduledAt
	}

	now := time.Now()

	return &Hearing{
		ID:          uuid.New(),
		CaseID:      caseID,
		Title:       title,
		Description: description,
		Type:        hearingType,
		Location:    location,
		ScheduledAt: scheduledAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
