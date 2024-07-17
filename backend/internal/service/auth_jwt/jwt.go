package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var Secret = []byte("secret") // secret key need to be stored somewhere

func generateAccessToken(id int, username string) (string, error) {
	accessClaims := jwt.MapClaims{
		"id":       id,
		"username": username,
		"exp":      time.Now().Add(2 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	at, err := accessToken.SignedString(Secret)
	if err != nil {
		return "", err
	}
	return at, nil
}

func generateRefreshToken(id int) (string, error) {
	refreshClaims := jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(24 * time.Hour * 30).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	rt, err := refreshToken.SignedString(Secret)
	if err != nil {
		return "", err
	}
	return rt, nil
}
