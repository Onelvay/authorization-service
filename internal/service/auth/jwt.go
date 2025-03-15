package auth

import (
	"account-service/internal/domain/grant"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func (s *Service) generateJWT(secretKey string, userID string) (string, error) {
	claims := grant.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
