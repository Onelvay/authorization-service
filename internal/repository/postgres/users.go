package postgres

import (
	"context"
	"time"

	"account-service/internal/domain/users"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(ctx context.Context, data users.User) (id string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
			INSERT INTO users (email, name, password)
			VALUES ($1, $2,$3)
			RETURNING id;`

	args := []interface{}{data.Email, data.Name, data.Password}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id)
	return
}

func (r *Repository) GetUsers(ctx context.Context) (dest []users.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
			SELECT created_at, updated_at, id, email, name, password
				FROM users`
	args := []interface{}{}

	err = r.db.SelectContext(ctx, &dest, query, args...)
	return
}

func (r *Repository) GetUserByEmailOrLogin(ctx context.Context, email string, login string) (dest users.User, err error) {
	query := `
			SELECT created_at, updated_at, id, email, name, password
			FROM users
			WHERE email = $1 OR name = $2;`

	args := []interface{}{email, login}

	err = r.db.GetContext(ctx, &dest, query, args...)
	return
}

func (r *Repository) GetUserByAny(ctx context.Context, login string) (dest users.User, err error) {
	query := `
			SELECT created_at, updated_at, id, email, name, password
			FROM users
			WHERE email = $1 OR name = $1;`

	args := []interface{}{login}

	err = r.db.GetContext(ctx, &dest, query, args...)
	return
}
