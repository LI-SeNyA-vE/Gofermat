package middleware

import (
	"fmt"
	"github.com/LI-SeNyA-vE/Gofermat/internal/service"
	"net/http"
	"strings"
)

func VerificationJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && !strings.HasPrefix(authHeader, "Bearer ") {
			jwtToken := strings.TrimPrefix(authHeader, "Bearer ")
			_, err := service.VerificationToken(jwtToken)
			if err != nil {
				http.Error(w, fmt.Sprintf("невалидный токен, ошибка: %s", err), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "не передан токен", http.StatusUnauthorized)
		return
	})
}
