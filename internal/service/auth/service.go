package auth

import (
	"account-service/internal/domain/users"
)

type Configuration func(s *Service) error

type Service struct {
	userRepository users.Repository
}

func New(configs ...Configuration) (s *Service, err error) {
	s = &Service{}

	for _, cfg := range configs {
		if err = cfg(s); err != nil {
			return
		}
	}

	return
}

func WithUserRepository(r users.Repository) Configuration {
	return func(s *Service) error {
		s.userRepository = r
		return nil
	}
}
