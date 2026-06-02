package document

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	caseDomain "legalflow/internal/domain/legalcase"
	documentDomain "legalflow/internal/domain/document"
	timelineDomain "legalflow/internal/domain/timeline"
)

type UploadDocumentUseCase struct {
	docRepo     documentDomain.Repository
	caseRepo    caseDomain.CaseRepository
	timelineRepo timelineDomain.Repository
}

func NewUploadDocumentUseCase(
	docRepo documentDomain.Repository,
	caseRepo caseDomain.CaseRepository,
	timelineRepo timelineDomain.Repository,
) *UploadDocumentUseCase {
	return &UploadDocumentUseCase{
		docRepo:      docRepo,
		caseRepo:     caseRepo,
		timelineRepo: timelineRepo,
	}
}

func (uc *UploadDocumentUseCase) Execute(ctx context.Context, input UploadDocumentInput) (*UploadDocumentOutput, error) {
	caseID, err := uuid.Parse(input.CaseID)
	if err != nil {
		return nil, ErrInvalidCaseID
	}

	c, err := uc.caseRepo.FindByID(ctx, caseID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrCaseNotFound
	}

	doc, err := documentDomain.NewDocument(caseID, input.Name, input.Description, documentDomain.DocumentType(input.Type), input.FileName, input.MimeType, input.FileSize, input.StorageKey)
	if err != nil {
		return nil, ErrInvalidInput
	}

	if err := uc.docRepo.Create(ctx, doc); err != nil {
		return nil, err
	}

	meta := timelineDomain.DocumentMetadata{
		DocumentID:   doc.ID.String(),
		DocumentType: string(doc.Type),
		FileName:     doc.FileName,
	}
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}

	event := timelineDomain.NewEventWithMetadata(caseID, timelineDomain.EventDocumentUploaded, fmt.Sprintf("Documento enviado: %s", input.Name), metaJSON)
	if err := uc.timelineRepo.Create(ctx, event); err != nil {
		return nil, err
	}

	return &UploadDocumentOutput{
		ID:          doc.ID.String(),
		CaseID:      doc.CaseID.String(),
		Name:        doc.Name,
		Description: doc.Description,
		Type:        string(doc.Type),
		FileName:    doc.FileName,
		MimeType:    doc.MimeType,
		FileSize:    doc.FileSize,
		StorageKey:  doc.StorageKey,
	}, nil
}
