package http

import (
	"account-service/internal/domain/grant"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"strings"
)

func extractUserIDFromJWT(r *http.Request) (string, error) {
	tokenString, err := ExtractBearerToken(r)
	if err != nil {
		return "", err
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &grant.Claims{}, keyFunc)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*grant.Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", errors.New("invalid token or claims")
}

func ExtractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", errors.New("invalid authorization header format")
	}

	return strings.TrimPrefix(authHeader, bearerPrefix), nil
}
