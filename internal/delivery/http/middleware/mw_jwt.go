package middleware

import (
	"fmt"
	"github.com/LI-SeNyA-vE/Gofermat/internal/global"
	"github.com/LI-SeNyA-vE/Gofermat/internal/service"
	"net/http"
	"strings"
)

func VerificationJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		global.Logger.Infof("Authorization = %s", authHeader)
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			jwtToken := strings.TrimPrefix(authHeader, "Bearer ")
			global.Logger.Infof("токен = %s", authHeader)
			_, err := service.VerificationToken(jwtToken)
			if err != nil {
				global.Logger.Infof("невалидный токен, ошибка: %s", err)
				http.Error(w, fmt.Sprintf("невалидный токен, ошибка: %s", err), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		}
		global.Logger.Infof("не передан токен в мидлваре")
		http.Error(w, "не передан токен", http.StatusUnauthorized)
		return
	})
}
