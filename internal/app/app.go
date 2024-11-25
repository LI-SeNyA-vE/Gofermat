package app

import (
	"github.com/LI-SeNyA-vE/Gofermat/internal/config"
	"github.com/LI-SeNyA-vE/Gofermat/internal/delivery/http/router"
	"github.com/LI-SeNyA-vE/Gofermat/internal/global"
	"github.com/LI-SeNyA-vE/Gofermat/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Run(configPath string) {
	//инициализация Логера
	config.Start(configPath)
	gormDB, err := service.Start()
	defer func() {
		sqlDB, _ := gormDB.DB()
		_ = sqlDB.Close()
	}()

	if err != nil {
		global.Logger.Infof("ошибка на моменте инициализации базы даннх\n%v", err)
	}

	//Создаёт роутер
	r := router.SetapRouter()

	//Старт сервера
	startServer(r)
}

func startServer(r *chi.Mux) {
	global.Logger.Infof("Открыт сервер %s", global.Config.Flags.RunAddress)
	err := http.ListenAndServe(global.Config.Flags.RunAddress, r)
	if err != nil {
		panic(err)
	}

}
