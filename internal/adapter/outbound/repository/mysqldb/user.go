package mysqldb

import (
	"context"
	"database/sql"
	"errors"

	"skeleton-service/internal/core/domain"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
}

type userRepository struct {
	db *sql.DB
}

func (r *repoSQL) GetUserMysqlDbRepository() UserRepository {
	return &userRepository{db: r.db}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT id, full_name, email, password_hash, created_at
		  FROM users
		 WHERE email = ?
		 LIMIT 1
	`

	var user domain.User
	if err := r.db.QueryRowContext(ctx, q, email).
		Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	const q = `
		INSERT INTO users (full_name, email, password_hash)
		VALUES (?, ?, ?)
	`

	res, err := r.db.ExecContext(ctx, q, user.FullName, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	const q = `
		SELECT id, full_name, email, password_hash, created_at
		  FROM users
		 WHERE id = ?
		 LIMIT 1
	`

	var user domain.User
	if err := r.db.QueryRowContext(ctx, q, id).
		Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
