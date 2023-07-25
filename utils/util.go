package utils

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"gorm-crud/config"
	"gorm-crud/models"
	"time"
)

var secretKey = []byte(config.Config("SECRET"))

func CreateToken(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	hashedPassword := fmt.Sprintf("%x", hash)
	return hashedPassword
}

func VerifyPassword(password, hashedPassword string) bool {
	hashedInput := HashPassword(password)
	if hashedInput == hashedPassword {
		return true
	}
	return false
}
