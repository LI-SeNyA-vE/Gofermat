package service

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const secretKey string = "test_secret_key"

type Claims struct {
	SecretKey string
}

func (c *Claims) CreateNewToken(email string) (string, error) {
	payload := jwt.MapClaims{
		"sub": email,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString([]byte(c.SecretKey))
	if err != nil {
		return "", err
	}
	return t, nil
}

func main() {
	token := Claims{SecretKey: secretKey}
	newToken, err := token.CreateNewToken("test@test.ru")
	if err != nil {
		println(err)
		return
	}
	println(newToken)
}
