// myblog-project/auth-service/internal/utils.go
package internal

import (
	"time"
	"os"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(username string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}