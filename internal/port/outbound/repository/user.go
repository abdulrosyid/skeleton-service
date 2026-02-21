package repository

import (
	"context"

	"skeleton-service/internal/core/domain"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
}
