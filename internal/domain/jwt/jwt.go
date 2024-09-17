package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func NewTokens(guid, ownKey string, ip string, expHours int) (string, string, string, error) {
	const fn = "domain.jwt.NewTokens"

	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)

	claims["guid"] = guid
	claims["ip"] = ip
	claims["exp"] = time.Now().Add(time.Duration(expHours) * time.Hour).Unix()

	tokenString, err := token.SignedString([]byte(ownKey))
	if err != nil {
		return "", "", "", fmt.Errorf("%s: %w", fn, err)
	}

	refreshToken := uuid.New().String()
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", "", fmt.Errorf("%s: %w", fn, err)
	}

	return tokenString, refreshToken, string(hashedRefreshToken), nil
}
