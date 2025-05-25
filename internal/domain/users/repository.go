package users

import (
	"account-service/internal/domain/billing"
	"context"
)

type Repository interface {
	CreateUser(ctx context.Context, data User) (string, error)
	GetUsers(ctx context.Context) ([]User, error)
	GetUserByAny(ctx context.Context, login string) (dest User, err error)
	GetUserByEmailOrLogin(ctx context.Context, email string, login string) (dest User, err error)

	UpdateUser(ctx context.Context, user User) error
	CreateSub(ctx context.Context, accId string, name string) (id string, err error)

	CreateBilling(ctx context.Context, data billing.Entity) (id string, err error)
	GetBillingByID(ctx context.Context, id string) (billing.Entity, error)
	CreateCard(ctx context.Context, data billing.CardEntity) (id string, err error)
	GetCards(ctx context.Context, accID string) (dest []billing.CardEntity, err error)
	DeleteCardByID(ctx context.Context, id string) error
}
