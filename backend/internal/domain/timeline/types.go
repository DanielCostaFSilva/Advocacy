package timeline

type EventType string

const (
	EventCaseCreated   EventType = "CASE_CREATED"
	EventStatusChanged EventType = "STATUS_CHANGED"
	EventNoteAdded     EventType = "NOTE_ADDED"
)
