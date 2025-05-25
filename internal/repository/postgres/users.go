package postgres

import (
	"context"
	"time"

	"account-service/internal/domain/billing"
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
			SELECT created_at, updated_at, id, email, name, password, phone,birth_date,gender
				FROM users`
	args := []interface{}{}

	err = r.db.SelectContext(ctx, &dest, query, args...)
	return
}

func (r *Repository) GetUserByEmailOrLogin(ctx context.Context, email string, login string) (dest users.User, err error) {
	query := `
			SELECT created_at, updated_at, id, email, name, password, phone,birth_date,gender
			FROM users
			WHERE email = $1 OR name = $2;`

	args := []interface{}{email, login}

	err = r.db.GetContext(ctx, &dest, query, args...)
	return
}

func (r *Repository) GetUserByAny(ctx context.Context, login string) (dest users.User, err error) {
	query := `
			SELECT created_at, updated_at, id, email, name, password, phone,birth_date,gender
			FROM users
			WHERE email = $1 OR name = $1;`

	args := []interface{}{login}

	err = r.db.GetContext(ctx, &dest, query, args...)
	return
}

func (r *Repository) CreateBilling(ctx context.Context, data billing.Entity) (id string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		INSERT INTO billings (
			iin, correlation_id, invoice_id, amount, currency,
			terminal_id, description, account_id, name, email,
			phone, back_link, failure_back_link, post_link, failure_post_link,
			language, data, card_save, source
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15,
			$16, $17, $18, $19
		)
		RETURNING id;
	`

	args := []interface{}{
		data.IIN, data.CorrelationID, data.InvoiceID, data.Amount, data.Currency,
		data.TerminalID, data.Description, data.AccountID, data.Name, data.Email,
		data.Phone, data.BackLink, data.FailureBackLink, data.PostLink, data.FailurePostLink,
		data.Language, data.Data, data.CardSave, data.Source,
	}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id)
	return
}

func (r *Repository) GetBillingByID(ctx context.Context, id string) (billing.Entity, error) {
	var b billing.Entity

	query := `
		SELECT 
			id, correlation_id, invoice_id, iin, phone, source,
			amount, currency, terminal_id, description, account_id,
			name, email, data, back_link, failure_back_link,
			post_link, failure_post_link, language, card_save
		FROM billings
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := r.db.GetContext(ctx, &b, query, id)
	return b, err
}

func (r *Repository) CreateCard(ctx context.Context, data billing.CardEntity) (id string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		INSERT INTO cards (
			card_id, account_id, terminal_id, type, mask, 
			issuer, is_default
		) VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7
		)
		RETURNING id;
	`

	args := []interface{}{
		data.CardID, data.AccountID, data.TerminalID, data.Type, data.Mask,
		data.Issuer, data.IsDefault,
	}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id)
	return
}

func (r *Repository) GetCards(ctx context.Context, accID string) (dest []billing.CardEntity, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
			SELECT id, account_id, terminal_id, type, mask, 
			issuer, is_default
				FROM cards WHERE account_id = $1;`
	args := []interface{}{accID}

	err = r.db.SelectContext(ctx, &dest, query, args...)
	return
}

func (r *Repository) DeleteCardByID(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		DELETE FROM cards
		WHERE id = $1;
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
