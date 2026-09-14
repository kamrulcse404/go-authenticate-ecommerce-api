package user

import (
	"context"
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("user not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) FindByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT
			id,
			name,
			email,
			role,
			status,
			email_verified,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	var u User

	err := r.db.QueryRowContext(
		ctx, query, id,
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}
