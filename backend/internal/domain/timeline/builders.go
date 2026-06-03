package timeline

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	contractDomain "legalflow/internal/domain/contract"
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

func NewContractUpdatedEvent(contract *contractDomain.Contract, oldAmount string) *TimelineEvent {
	meta := ContractUpdatedMetadata{
		ContractID:   contract.ID.String(),
		OldAmount:    oldAmount,
		NewAmount:    contract.Amount.String(),
		ContractType: string(contract.Type),
	}
	metaJSON, _ := json.Marshal(meta)
	return NewEventWithMetadata(contract.CaseID, EventContractUpdated, fmt.Sprintf("Contrato atualizado: %s", contract.Title), metaJSON)
}

func NewContractClosedEvent(contract *contractDomain.Contract) *TimelineEvent {
	meta := ContractClosedMetadata{
		ContractID:   contract.ID.String(),
		ContractType: string(contract.Type),
		Amount:       contract.Amount.String(),
		ClosedAt:     time.Now().UTC().Format(time.RFC3339),
	}
	metaJSON, _ := json.Marshal(meta)
	return NewEventWithMetadata(contract.CaseID, EventContractClosed, fmt.Sprintf("Contrato encerrado: %s", contract.Title), metaJSON)
}

func NewContractCreatedEvent(caseID uuid.UUID, contract *contractDomain.Contract) *TimelineEvent {
	meta := ContractMetadata{
		ContractID:   contract.ID.String(),
		ContractType: string(contract.Type),
		Amount:       contract.Amount.String(),
	}
	metaJSON, _ := json.Marshal(meta)
	return NewEventWithMetadata(caseID, EventContractCreated, fmt.Sprintf("Contrato criado: %s", contract.Title), metaJSON)
}
