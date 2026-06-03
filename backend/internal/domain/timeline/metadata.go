package timeline

import "time"

type HearingMetadata struct {
	HearingID   string    `json:"hearing_id"`
	Type        string    `json:"hearing_type"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

type DocumentMetadata struct {
	DocumentID   string `json:"document_id"`
	DocumentType string `json:"document_type"`
	FileName     string `json:"file_name"`
}

type DocumentDownloadMetadata struct {
	DocumentID string `json:"document_id"`
	FileName   string `json:"file_name"`
}

type ContractMetadata struct {
	ContractID   string `json:"contract_id"`
	ContractType string `json:"contract_type"`
	Amount       string `json:"amount"`
}

type ContractUpdatedMetadata struct {
	ContractID   string `json:"contract_id"`
	OldAmount    string `json:"old_amount"`
	NewAmount    string `json:"new_amount"`
	ContractType string `json:"contract_type"`
}

type ContractClosedMetadata struct {
	ContractID   string `json:"contract_id"`
	ContractType string `json:"contract_type"`
	Amount       string `json:"amount"`
	ClosedAt     string `json:"closed_at"`
}
