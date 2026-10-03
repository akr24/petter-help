package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/akr24/petter-help/backend/internal/domain"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

const userColumns = `id, email, name, password_hash, role, created_at`

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO users (email, name, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`,
		u.Email, u.Name, u.PasswordHash, u.Role,
	).Scan(&u.ID, &u.CreatedAt)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
		return domain.ErrEmailTaken
	}
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	return r.findOne(ctx, `SELECT `+userColumns+` FROM users WHERE id = $1`, id)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.findOne(ctx, `SELECT `+userColumns+` FROM users WHERE email = $1`, email)
}

func (r *UserRepository) findOne(ctx context.Context, sql string, arg any) (*domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx, sql, arg).
		Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
