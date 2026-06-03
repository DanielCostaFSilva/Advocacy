package document

import (
	"context"
	"time"

	"github.com/google/uuid"
	caseDomain "legalflow/internal/domain/legalcase"
	documentDomain "legalflow/internal/domain/document"
)

type ListCaseDocumentsInput struct {
	CaseID   string
	Page     int
	PageSize int
	Type     string
	Sort     string
	Order    string
}

type DocumentDTO struct {
	ID         string    `json:"id"`
	CaseID     string    `json:"case_id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	FileName   string    `json:"file_name"`
	MimeType   string    `json:"mime_type"`
	FileSize   int64     `json:"file_size"`
	StorageKey string    `json:"storage_key"`
	CreatedAt  time.Time `json:"created_at"`
}

type ListCaseDocumentsOutput struct {
	Data       []DocumentDTO
	Page       int
	PageSize   int
	Total      int64
	TotalPages int
}

type ListCaseDocumentsUseCase struct {
	docRepo  documentDomain.Repository
	caseRepo caseDomain.CaseRepository
}

func NewListCaseDocumentsUseCase(docRepo documentDomain.Repository, caseRepo caseDomain.CaseRepository) *ListCaseDocumentsUseCase {
	return &ListCaseDocumentsUseCase{
		docRepo:  docRepo,
		caseRepo: caseRepo,
	}
}

func (uc *ListCaseDocumentsUseCase) Execute(ctx context.Context, input ListCaseDocumentsInput) (*ListCaseDocumentsOutput, error) {
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

	page := input.Page
	if page < 1 {
		page = 1
	}

	pageSize := input.PageSize
	if pageSize < 1 {
		pageSize = 20
	} else if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	params := documentDomain.ListDocumentsParams{
		Offset: offset,
		Limit:  pageSize,
		CaseID: input.CaseID,
		Type:   input.Type,
		Sort:   input.Sort,
		Order:  input.Order,
	}

	documents, total, err := uc.docRepo.ListByCaseIDPaginated(ctx, params)
	if err != nil {
		return nil, err
	}

	totalPages := int(total / int64(pageSize))
	if total%int64(pageSize) > 0 {
		totalPages++
	}

	data := make([]DocumentDTO, len(documents))
	for i, d := range documents {
		data[i] = DocumentDTO{
			ID:         d.ID.String(),
			CaseID:     d.CaseID.String(),
			Name:       d.Name,
			Type:       string(d.Type),
			FileName:   d.FileName,
			MimeType:   d.MimeType,
			FileSize:   d.FileSize,
			StorageKey: d.StorageKey,
			CreatedAt:  d.CreatedAt,
		}
	}

	return &ListCaseDocumentsOutput{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
