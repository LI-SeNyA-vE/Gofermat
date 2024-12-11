package app

import (
	"github.com/LI-SeNyA-vE/Gofermat/internal/config"
	"github.com/LI-SeNyA-vE/Gofermat/internal/delivery/http/router"
	"github.com/LI-SeNyA-vE/Gofermat/internal/global"
	"github.com/LI-SeNyA-vE/Gofermat/internal/model"
	"github.com/LI-SeNyA-vE/Gofermat/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Run() {
	//инициализация Логера
	configs, err := config.Start()

	if err != nil {
		global.Logger.Infof("ошибка на моменте загрузки конфига\n%v", err)
	}

	gormDB, err := service.Start()
	defer func() {
		sqlDB, _ := gormDB.DB()
		_ = sqlDB.Close()
	}()

	if err != nil {
		global.Logger.Infof("ошибка на моменте инициализации базы данных\n%v", err)
	}

	//Создаёт роутер
	r := router.SetupRouter()

	//Старт сервера
	startServer(r, configs.ConfigFlag)
}

func startServer(r *chi.Mux, conf model.ConfigFlag) {
	global.Logger.Infof("Открыт сервер %s", conf.RunAddress)
	err := http.ListenAndServe(conf.RunAddress, r)
	if err != nil {
		panic(err)
	}
}
