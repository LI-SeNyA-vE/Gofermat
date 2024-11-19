package service

import (
	"errors"
	"fmt"
	"github.com/LI-SeNyA-vE/Gofermat/internal/global"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func CreateNewToken(login string) (string, error) {
	payload := jwt.MapClaims{
		"sub": login,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString([]byte(global.Config.SecretKey))
	if err != nil {
		return "", fmt.Errorf("ошибка: %w при создании токена для пользователя: %s", err, login)
	}
	return t, nil
}

func VerificationToken(userToken string) (jwt.MapClaims, error) {

	t, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.SecretKey), nil
	})

	switch {
	case t.Valid:
		global.Logger.Info("токен правильный")
		claims := t.Claims.(jwt.MapClaims)
		return claims, nil
	case errors.Is(err, jwt.ErrTokenMalformed):
		return nil, fmt.Errorf("токен имеет неправильную форму %w", err)
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return nil, fmt.Errorf("подпись токена недействительна %w", err)
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return nil, fmt.Errorf("срок действия токена истек или токен еще не действителен %w", err)
	default:
		return nil, fmt.Errorf("неизвестная ошика на этапе проверки валидности токена: %w", err)
	}
}
