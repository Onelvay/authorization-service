package repository

import (
	"account-service/internal/domain/users"
	"account-service/internal/repository/postgres"
	"account-service/pkg/store"
)

type Configuration func(r *Repository) error

type Repository struct {
	postgres store.SQL

	User users.Repository
}

func New(configs ...Configuration) (s *Repository, err error) {
	s = &Repository{}

	// Apply all Configurations passed in
	for _, cfg := range configs {
		if err = cfg(s); err != nil {
			return
		}
	}

	return
}

func (r *Repository) Close() {
	if r.postgres.Connection != nil {
		r.postgres.Connection.Close()
	}
}

func WithPostgresStore(dataSourceName string) Configuration {
	return func(r *Repository) (err error) {
		r.postgres, err = store.NewSQL(dataSourceName)
		if err != nil {
			return
		}

		if err = store.Migrate(dataSourceName); err != nil {
			return
		}

		r.User = postgres.NewUserRepository(r.postgres.Connection)
		return
	}
}
