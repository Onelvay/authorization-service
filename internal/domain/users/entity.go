package users

import (
	"account-service/internal/domain/grant"
	"time"
)

type User struct {
	ID        string    `db:"id,primarykey"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	Password  string    `db:"password"`
	Phone     *string   `db:"phone"`
	BirthDate *string   `db:"birth_date"`
	Gender    *string   `db:"gender"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func ParseFromAuth(req grant.Request) User {
	return User{
		Name:     req.Login,
		Email:    req.Email,
		Password: req.Password,
	}
}
