package router

import (
	"github.com/LI-SeNyA-vE/Gofermat/internal/delivery/http/handlers"
	"github.com/LI-SeNyA-vE/Gofermat/internal/delivery/http/middleware"
	"github.com/go-chi/chi/v5"
)

//POST /api/user/register — регистрация пользователя;
//POST /api/user/login — аутентификация пользователя;
//POST /api/user/orders — загрузка пользователем номера заказа для расчёта;
//GET /api/user/orders — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
//GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
//POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
//GET /api/user/withdrawals — получение информации о выводе средств с накопительного счёта пользователем.

func SetapRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(middleware.LoggingMiddleware)
		r.Use(middleware.GzipMiddleware)
		r.Use(middleware.UnGzipMiddleware)

		r.Post("/api/user/register", handlers.UserRegistration) //регистрация пользователя
		r.Post("/api/user/login", handlers.UserAuthentication)  //аутентификация пользователя

		r.Group(func(r chi.Router) {
			r.Use(middleware.VerificationJWT)

			r.Post("/api/user/orders", handlers.AddOrder)                              //загрузка пользователем номера заказа для расчёта
			r.Post("/api/user/balance/withdraw", handlers.ExpenditurePointsOnNewOrder) //запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
			r.Get("/api/user/orders", handlers.ListUserOrders)                         //получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
			r.Get("/api/user/balance", handlers.ListUserBalance)                       //получение текущего баланса счёта баллов лояльности пользователя
			r.Get("/api/user/withdrawals", handlers.InfoAboutUsagePoints)              //получение информации о выводе средств с накопительного счёта пользователем
		})

	})
	return r
}
