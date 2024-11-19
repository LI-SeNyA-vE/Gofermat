package middleware

import (
	"fmt"
	"github.com/LI-SeNyA-vE/Gofermat/internal/service"
	"net/http"
)

func VerificationJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken := r.Header.Get("Authorization")
		if jwtToken != "" {
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
