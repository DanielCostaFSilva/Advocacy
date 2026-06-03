package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	domain "legalflow/internal/domain/contract"
)

type ContractRepository struct {
	db *sql.DB
	q  *Queries
}

func NewContractRepository(db *sql.DB) *ContractRepository {
	return &ContractRepository{
		db: db,
		q:  New(db),
	}
}

func (r *ContractRepository) Create(ctx context.Context, contract *domain.Contract) error {
	var description sql.NullString
	if contract.Description != "" {
		description = sql.NullString{String: contract.Description, Valid: true}
	}

	var endDate sql.NullTime
	if contract.EndDate != nil {
		endDate = sql.NullTime{Time: *contract.EndDate, Valid: true}
	}

	created, err := r.q.CreateContract(ctx, CreateContractParams{
		ClientID:    contract.ClientID,
		CaseID:      contract.CaseID,
		Title:       contract.Title,
		Description: description,
		Type:        string(contract.Type),
		Amount:      contract.Amount.String(),
		StartDate:   contract.StartDate,
		EndDate:     endDate,
	})
	if err != nil {
		return err
	}

	contract.ID = created.ID
	contract.CreatedAt = created.CreatedAt
	contract.UpdatedAt = created.UpdatedAt
	return nil
}

func (r *ContractRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Contract, error) {
	result, err := r.q.GetContractByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toContractDomain(result), nil
}

func (r *ContractRepository) ListByCaseID(ctx context.Context, caseID uuid.UUID) ([]domain.Contract, error) {
	rows, err := r.q.ListContractsByCaseID(ctx, caseID)
	if err != nil {
		return nil, err
	}
	contracts := make([]domain.Contract, len(rows))
	for i, row := range rows {
		contracts[i] = *toContractDomain(row)
	}
	return contracts, nil
}

func toContractDomain(c Contract) *domain.Contract {
	amount, _ := decimal.NewFromString(c.Amount)

	var endDate *time.Time
	if c.EndDate.Valid {
		endDate = &c.EndDate.Time
	}

	return &domain.Contract{
		ID:          c.ID,
		ClientID:    c.ClientID,
		CaseID:      c.CaseID,
		Title:       c.Title,
		Description: c.Description.String,
		Type:        domain.ContractType(c.Type),
		Amount:      amount,
		StartDate:   c.StartDate,
		EndDate:     endDate,
		Active:      c.Active,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
