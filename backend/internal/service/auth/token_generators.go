package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kotopesp/sos-kotopes/pkg/logger"

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

func (s *service) generateHash(str string) string {
	hashedToken := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hashedToken[:])
}

// generateRefreshToken Generating refresh token
func (s *service) generateRefreshToken() string {
	return uuid.New().String()
}
