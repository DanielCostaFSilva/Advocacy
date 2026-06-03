package document

import (
	"context"
	"time"

	"github.com/google/uuid"
	documentDomain "legalflow/internal/domain/document"
	timelineDomain "legalflow/internal/domain/timeline"
)

type GetDocumentDownloadInput struct {
	DocumentID string
}

type GetDocumentDownloadOutput struct {
	DocumentID  string
	FileName    string
	DownloadURL string
	ExpiresAt   time.Time
}

type GetDocumentDownloadUseCase struct {
	docRepo     documentDomain.Repository
	timelineRepo timelineDomain.Repository
	urlProvider  DownloadURLProvider
}

func NewGetDocumentDownloadUseCase(
	docRepo documentDomain.Repository,
	timelineRepo timelineDomain.Repository,
	urlProvider DownloadURLProvider,
) *GetDocumentDownloadUseCase {
	return &GetDocumentDownloadUseCase{
		docRepo:      docRepo,
		timelineRepo: timelineRepo,
		urlProvider:  urlProvider,
	}
}

func (uc *GetDocumentDownloadUseCase) Execute(ctx context.Context, input GetDocumentDownloadInput) (*GetDocumentDownloadOutput, error) {
	docID, err := uuid.Parse(input.DocumentID)
	if err != nil {
		return nil, ErrInvalidDocumentID
	}

	doc, err := uc.docRepo.FindByID(ctx, docID)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, ErrDocumentNotFound
	}

	downloadURL := uc.urlProvider.Generate(doc)
	expiresAt := uc.urlProvider.ExpiresAt()

	event := timelineDomain.NewDocumentDownloadedEvent(doc.CaseID, doc)
	if err := uc.timelineRepo.Create(ctx, event); err != nil {
		return nil, ErrTimelinePersistence
	}

	return &GetDocumentDownloadOutput{
		DocumentID:  doc.ID.String(),
		FileName:    doc.FileName,
		DownloadURL: downloadURL,
		ExpiresAt:   expiresAt,
	}, nil
}
