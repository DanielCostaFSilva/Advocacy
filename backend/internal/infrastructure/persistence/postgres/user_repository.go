package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	domain "legalflow/internal/user/domain/user"
)

type UserRepository struct {
	q *Queries
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		q: New(db),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	params := CreateUserParams{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}

	created, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	user.ID = created.ID
	user.CreatedAt = created.CreatedAt
	user.UpdatedAt = created.UpdatedAt

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	result, err := r.q.FindUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return toDomain(result), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	result, err := r.q.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return toDomain(result), nil
}

func toDomain(u User) *domain.User {
	return &domain.User{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
