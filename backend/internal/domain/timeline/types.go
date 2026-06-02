package timeline

type EventType string

const (
	EventCaseCreated        EventType = "CASE_CREATED"
	EventStatusChanged      EventType = "STATUS_CHANGED"
	EventNoteAdded          EventType = "NOTE_ADDED"
	EventHearingCreated     EventType = "HEARING_CREATED"
	EventHearingRescheduled EventType = "HEARING_RESCHEDULED"
	EventHearingCancelled   EventType = "HEARING_CANCELLED"
	EventDocumentUploaded   EventType = "DOCUMENT_UPLOADED"
)
