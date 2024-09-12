package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"time"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func GenerateJWT(guid, ownKey string, ip string, expHours int) (string, error) {
	const fn = "domain.jwt.GenerateJWT"
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)
	claims["guid"] = guid
	claims["ip"] = ip
	claims["exp"] = time.Now().Add(time.Duration(expHours) * time.Hour).Unix()

	signedToken, err := token.SignedString([]byte(ownKey))
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}
	return signedToken, nil
}

func GenerateRefreshToken() (string, string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", "", err
	}

	refreshToken := base64.StdEncoding.EncodeToString(tokenBytes)

	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return refreshToken, string(hashedRefreshToken), nil
}
