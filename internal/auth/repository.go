package auth

import (
	"context"
	"database/sql"
	"ecommerce-api/internal/user"
	"errors"

	"github.com/lib/pq"
)

var ErrEmailAlreadyExists = errors.New("email already exists")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(ctx context.Context, name string, email string, passwordHash string) (*user.User, error) {
	query := `
		INSERT INTO users (
			name,
			email,
			password_hash
		)
		VALUES($1, $2, $3)
		RETURNING
			id,
			name,
			email,
			role,
			status,
			email_verified,
			created_at,
			updated_at
	`
	var u user.User

	err := r.db.QueryRowContext(
		ctx,
		query,
		name,
		email,
		passwordHash,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Role,
		&u.Status,
		&u.EmailVerified,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && (pqErr.Code == "23505") {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	return &u, nil
}
