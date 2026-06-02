package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

func (r *DocumentRepository) ListByCaseIDPaginated(ctx context.Context, params domain.ListDocumentsParams) ([]domain.Document, int64, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if params.CaseID != "" {
		conditions = append(conditions, fmt.Sprintf("case_id = $%d", argIdx))
		args = append(args, params.CaseID)
		argIdx++
	}
	if params.Type != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIdx))
		args = append(args, params.Type)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	sort := "created_at"
	switch params.Sort {
	case "name":
		sort = "name"
	case "type":
		sort = "type"
	case "created_at":
		sort = "created_at"
	case "file_size":
		sort = "file_size"
	}

	order := "asc"
	if params.Order == "desc" {
		order = "desc"
	}

	query := fmt.Sprintf(
		"SELECT id, case_id, name, description, type, file_name, mime_type, file_size, storage_key, created_at, updated_at FROM documents %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		whereClause, sort, order, argIdx, argIdx+1,
	)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var documents []domain.Document
	for rows.Next() {
		var d Document
		if err := rows.Scan(&d.ID, &d.CaseID, &d.Name, &d.Description, &d.Type, &d.FileName, &d.MimeType, &d.FileSize, &d.StorageKey, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, 0, err
		}
		documents = append(documents, *toDocumentDomain(d))
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM documents %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total); err != nil {
		return nil, 0, err
	}

	return documents, total, nil
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
