package types

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenType string

const (
	AccessToken  TokenType = "access_token"
	RefreshToken TokenType = "refresh_token"
)

func GenerateToken(user *User, tokenType TokenType) (string, error) {
	var expirationTime time.Time

	switch tokenType {
	case AccessToken:
		expirationTime = time.Now().Add(time.Minute * 30)
	case RefreshToken:
		expirationTime = time.Now().Add(time.Hour * 24 * 7)
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expirationTime.Unix(),
		"type":    tokenType,
	}

	if tokenType == AccessToken {
		claims["role"] = user.Role
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var secret string
	if tokenType == AccessToken {
		secret = os.Getenv("JWT_SECRET")
	} else {
		secret = os.Getenv("JWT_REFRESH_SECRET")
	}

	return token.SignedString([]byte(secret))
}

func CreateTokenPair(user *User) (*TokenPair, error) {
	accessToken, err := GenerateToken(user, AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Error generating access token: %s", err)
	}

	refreshToken, err := GenerateToken(user, RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("Error generating refresh token: %s", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
