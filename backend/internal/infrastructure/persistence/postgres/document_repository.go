package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	domain "legalflow/internal/domain/document"
)

type DocumentRepository struct {
	db *sql.DB
	q  *Queries
}

func NewDocumentRepository(db *sql.DB) *DocumentRepository {
	return &DocumentRepository{
		db: db,
		q:  New(db),
	}
}

func (r *DocumentRepository) Create(ctx context.Context, d *domain.Document) error {
	params := CreateDocumentParams{
		CaseID:      d.CaseID,
		Name:        d.Name,
		Description: sql.NullString{String: d.Description, Valid: d.Description != ""},
		Type:        string(d.Type),
		FileName:    d.FileName,
		MimeType:    d.MimeType,
		FileSize:    d.FileSize,
		StorageKey:  d.StorageKey,
	}

	created, err := r.q.CreateDocument(ctx, params)
	if err != nil {
		return err
	}

	d.ID = created.ID
	d.CreatedAt = created.CreatedAt
	d.UpdatedAt = created.UpdatedAt
	return nil
}

func (r *DocumentRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Document, error) {
	result, err := r.q.GetDocumentByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return toDocumentDomain(result), nil
}

func (r *DocumentRepository) ListByCaseID(ctx context.Context, caseID uuid.UUID) ([]domain.Document, error) {
	rows, err := r.q.ListDocumentsByCaseID(ctx, caseID)
	if err != nil {
		return nil, err
	}

	docs := make([]domain.Document, len(rows))
	for i, row := range rows {
		docs[i] = *toDocumentDomain(row)
	}
	return docs, nil
}

func toDocumentDomain(d Document) *domain.Document {
	return &domain.Document{
		ID:          d.ID,
		CaseID:      d.CaseID,
		Name:        d.Name,
		Description: d.Description.String,
		Type:        domain.DocumentType(d.Type),
		FileName:    d.FileName,
		MimeType:    d.MimeType,
		FileSize:    d.FileSize,
		StorageKey:  d.StorageKey,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}
