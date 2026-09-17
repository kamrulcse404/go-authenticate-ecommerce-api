package auth

import (
	"context"
	"database/sql"
	"ecommerce-api/internal/user"
	"errors"

	"github.com/lib/pq"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserWithPassword struct {
	User         user.User
	PasswordHash string
}

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

func (r *Repository) FindByEmail(ctx context.Context, email string) (*UserWithPassword, error) {
	query := `
		SELECT 
			id,
			name,
			email,
			password_hash,
			role,
			status,
			email_verified,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	var result UserWithPassword

	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&result.User.ID,
		&result.User.Name,
		&result.User.Email,
		&result.PasswordHash,
		&result.User.Role,
		&result.User.Status,
		&result.User.EmailVerified,
		&result.User.CreatedAt,
		&result.User.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	return &result, nil
}
