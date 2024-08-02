package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *service) generateAccessToken(id int, username string) (*string, error) {
	accessClaims := jwt.MapClaims{
		"id":       id,
		"username": username,
		"exp":      time.Now().Add(2 * time.Minute).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	at, err := accessToken.SignedString(s.authServiceConfig.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &at, nil
}

func (s *service) generateRefreshToken(id int) (*string, error) {
	refreshClaims := jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(24 * time.Hour * 30).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	rt, err := refreshToken.SignedString(s.authServiceConfig.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &rt, nil
}
