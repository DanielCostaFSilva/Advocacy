package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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

func (r *ContractRepository) Update(ctx context.Context, contract *domain.Contract) error {
	var description sql.NullString
	if contract.Description != "" {
		description = sql.NullString{String: contract.Description, Valid: true}
	}

	var endDate sql.NullTime
	if contract.EndDate != nil {
		endDate = sql.NullTime{Time: *contract.EndDate, Valid: true}
	}

	updated, err := r.q.UpdateContract(ctx, UpdateContractParams{
		ID:          contract.ID,
		Title:       contract.Title,
		Description: description,
		Type:        string(contract.Type),
		Amount:      contract.Amount.String(),
		StartDate:   contract.StartDate,
		EndDate:     endDate,
		Active:      contract.Active,
	})
	if err != nil {
		return err
	}

	contract.UpdatedAt = updated.UpdatedAt
	return nil
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

func (r *ContractRepository) List(ctx context.Context, params domain.ListContractsParams) ([]domain.Contract, int64, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if params.ClientID != "" {
		conditions = append(conditions, fmt.Sprintf("client_id = $%d", argIdx))
		args = append(args, params.ClientID)
		argIdx++
	}
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
	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("active = $%d", argIdx))
		args = append(args, *params.Active)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	sort := "created_at"
	switch params.Sort {
	case "created_at":
		sort = "created_at"
	case "title":
		sort = "title"
	case "amount":
		sort = "amount"
	case "start_date":
		sort = "start_date"
	}

	order := "asc"
	if params.Order == "desc" {
		order = "desc"
	}

	query := fmt.Sprintf(
		"SELECT id, client_id, case_id, title, description, type, amount, start_date, end_date, active, created_at, updated_at FROM contracts %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		whereClause, sort, order, argIdx, argIdx+1,
	)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var contracts []domain.Contract
	for rows.Next() {
		var c Contract
		if err := rows.Scan(
			&c.ID, &c.ClientID, &c.CaseID, &c.Title, &c.Description,
			&c.Type, &c.Amount, &c.StartDate, &c.EndDate, &c.Active,
			&c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		contracts = append(contracts, *toContractDomain(c))
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM contracts %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total); err != nil {
		return nil, 0, err
	}

	return contracts, total, nil
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
