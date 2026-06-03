package timeline

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	documentDomain "legalflow/internal/domain/document"
)

func NewDocumentUploadedEvent(caseID uuid.UUID, document *documentDomain.Document) *TimelineEvent {
	meta := DocumentMetadata{
		DocumentID:   document.ID.String(),
		DocumentType: string(document.Type),
		FileName:     document.FileName,
	}
	metaJSON, _ := json.Marshal(meta)
	return NewEventWithMetadata(caseID, EventDocumentUploaded, fmt.Sprintf("Documento enviado: %s", document.Name), metaJSON)
}

func NewDocumentDownloadedEvent(caseID uuid.UUID, document *documentDomain.Document) *TimelineEvent {
	meta := DocumentMetadata{
		DocumentID:   document.ID.String(),
		DocumentType: string(document.Type),
		FileName:     document.FileName,
	}
	metaJSON, _ := json.Marshal(meta)
	return NewEventWithMetadata(caseID, EventDocumentDownloaded, fmt.Sprintf("Documento acessado: %s", document.Name), metaJSON)
}
