package auth

import (
	"account-service/internal/domain/grant"
	"account-service/internal/domain/users"
	"account-service/pkg/log"
	"context"
	"database/sql"
	"errors"
	"go.uber.org/zap"
	"os"
)

func (s *Service) SignUp(ctx context.Context, req grant.Request) (dest grant.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("SignUp")

	_, err = s.userRepository.GetUserByEmailOrLogin(ctx, req.Email, req.Login)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		data := users.ParseFromAuth(req)

		data.Password, err = hashPassword(data.Password)
		if err != nil {
			logger.Error("failed to hash", zap.Error(err))
			return
		}

		data.ID, err = s.userRepository.CreateUser(ctx, data)
		if err != nil {
			logger.Error("failed to create user", zap.Error(err))
			return
		}

		dest.AccessToken, err = s.generateJWT(os.Getenv("JWT_SECRET_KEY"), data.ID)
		if err != nil {
			logger.Error("failed to generate jwt", zap.Error(err))
			return
		}
	case errors.Is(err, nil):
		err = grant.ErrUserExist
	default:
		logger.Error("failed to get user", zap.Error(err))
	}

	return
}

func (s *Service) SignIn(ctx context.Context, req grant.Request) (dest grant.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("SignIn")

	data, err := s.userRepository.GetUserByAny(ctx, req.Login)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.Error("failed to get user", zap.Error(err))
		}
		return
	}

	if !checkPasswordHash(req.Password, data.Password) {
		err = errors.New("invalid password")
		return
	}

	dest.AccessToken, err = s.generateJWT(os.Getenv("JWT_SECRET_KEY"), data.ID)
	if err != nil {
		logger.Error("failed to generate jwt", zap.Error(err))
		return
	}

	return
}
