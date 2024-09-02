package auth

import (
	"context"
	"fmt"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// generateAccessToken Generating access token
func (s *service) generateAccessToken(id int, username string) (*string, error) {
	accessClaims := jwt.MapClaims{
		"id":       id,
		"username": username,
		"exp":      time.Now().Add(time.Duration(s.authServiceConfig.AccessTokenLifetime) * time.Minute).Unix(),
	}

	logger.Log().Debug(context.Background(), fmt.Sprintf("%d %d", time.Now().Unix(), accessClaims["exp"]))

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	at, err := accessToken.SignedString(s.authServiceConfig.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &at, nil
}

// generateRefreshToken Generating refresh token
func (s *service) generateRefreshToken(id int) (*string, error) {
	refreshClaims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Duration(s.authServiceConfig.RefreshTokenLifetime) * time.Minute).Unix(),
	}

	logger.Log().Debug(context.Background(), fmt.Sprintf("%d %d", time.Now().Unix(), refreshClaims["exp"]))

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	rt, err := refreshToken.SignedString(s.authServiceConfig.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &rt, nil
}
