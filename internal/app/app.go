package app

import (
	"github.com/LI-SeNyA-vE/Gofermat/internal/config"
	"github.com/LI-SeNyA-vE/Gofermat/internal/delivery/http/router"
	"github.com/LI-SeNyA-vE/Gofermat/internal/global"
	"github.com/LI-SeNyA-vE/Gofermat/internal/service"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func Run(configPath string) {
	server := Server{}
	server.Config = config.New()
	server.StorageDB = service.NewStorage()
	server.Config.LoadConfigFromFile("/Users/senya/GolandProjects/Gofermat/config.yml")
	err := server.StorageDB.Start()

	if err != nil {
		global.Logger.Infof("ошибка на моменте инициализации базы даннх\n%v", err)
	}

	//Создаёт роутер
	r := router.SetapRouter()

	//Старт сервера
	startServer(r)
}

type Server struct {
	StorageDB *service.Storage
	Config    *config.Config
}

//func NewServer(configPath string) (*Server, error) {
//	// Загружаем конфигурацию
//	config, err := configs.LoadConfigFromFile(configPath)
//	if err != nil {
//		return nil, err
//	}
//
//	// Создаём экземпляр Storage с загруженной конфигурацией
//	storageDB, err := storage.NewStorage(config)
//	if err != nil {
//		return nil, err
//	}
//
//	// Инициализируем базу данных
//	if err := storageDB.Start(); err != nil {
//		global.Logger.Infof("Ошибка на этапе инициализации базы данных: %v", err)
//		return nil, err
//	}
//
//	return &Server{
//		StorageDB: storageDB,
//	}, nil
//}

func startServer(r *chi.Mux) {
	log.Println("Открыт сервер ")
	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		panic(err)
	}

}
