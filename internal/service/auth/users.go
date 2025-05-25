package auth

import (
	"account-service/internal/domain/users"
	"account-service/pkg/log"
	"context"
	"go.uber.org/zap"
)

func (s *Service) GetUser(ctx context.Context, login string) (userData users.User, err error) {
	logger := log.LoggerFromContext(ctx).Named("get user")

	userData, err = s.userRepository.GetUserByEmailOrLogin(ctx, login, login)
	if err != nil {
		logger.Error("asas", zap.Error(err))
		return
	}

	return
}

func (s *Service) UpdateUser(ctx context.Context, userData users.User) (err error) {
	logger := log.LoggerFromContext(ctx).Named("update user")

	if err = s.userRepository.UpdateUser(ctx, userData); err != nil {
		logger.Error("asaszczxcz", zap.Error(err))
		return
	}
	return
}
