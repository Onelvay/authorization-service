package users

import "context"

type Repository interface {
	CreateUser(ctx context.Context, data User) (string, error)
	GetUsers(ctx context.Context) ([]User, error)
	GetUserByAny(ctx context.Context, login string) (dest User, err error)
	GetUserByEmailOrLogin(ctx context.Context, email string, login string) (dest User, err error)
}
