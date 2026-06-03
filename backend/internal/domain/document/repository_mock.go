package document

import (
	"context"

	"github.com/google/uuid"
)

type MockRepository struct {
	CreateFunc              func(ctx context.Context, document *Document) error
	FindByIDFunc            func(ctx context.Context, id uuid.UUID) (*Document, error)
	ListByCaseIDFunc        func(ctx context.Context, caseID uuid.UUID) ([]Document, error)
	ListByCaseIDPaginatedFunc func(ctx context.Context, params ListDocumentsParams) ([]Document, int64, error)
}

func (m *MockRepository) Create(ctx context.Context, document *Document) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, document)
	}
	return nil
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*Document, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockRepository) ListByCaseID(ctx context.Context, caseID uuid.UUID) ([]Document, error) {
	if m.ListByCaseIDFunc != nil {
		return m.ListByCaseIDFunc(ctx, caseID)
	}
	return nil, nil
}

func (m *MockRepository) ListByCaseIDPaginated(ctx context.Context, params ListDocumentsParams) ([]Document, int64, error) {
	if m.ListByCaseIDPaginatedFunc != nil {
		return m.ListByCaseIDPaginatedFunc(ctx, params)
	}
	return nil, 0, nil
}
