package config

import (
	"flag"
	"github.com/LI-SeNyA-vE/Gofermat/internal/global"
	"github.com/LI-SeNyA-vE/Gofermat/internal/logger"
	"gopkg.in/yaml.v3"
	"os"
)

func Start(fileName string) {
	err := logger.Initialize("info")
	if err != nil {
		return
	}
	LoadConfigFromFile(fileName)
	global.Config.Flags = loadConfig()
}

func LoadConfigFromFile(fileName string) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		global.Logger.Errorf("ошибка чтения конфигурационного файла %s\bОшибка: %v", fileName, err)
	}

	err = yaml.Unmarshal(data, &global.Config)
	if err != nil {
		global.Logger.Errorf("ошибка парсинга Yaml файла в config\b%v", err)
	}
}

func loadConfig() global.ConfigFlag {
	// Устанавливаем значения по умолчанию
	defaultRunAddress := "localhost:8080"
	defaultDatabaseURI := "postgres://senya:1q2w3e4r5t@localhost:5433/gofermat"
	defaultAccrualSystemAddress := "http://localhost:8081"

	// Читаем переменные окружения
	envRunAddress := os.Getenv("RUN_ADDRESS")
	envDatabaseURI := os.Getenv("DATABASE_URI")
	envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")

	// Определяем флаги
	flagRunAddress := flag.String("a", defaultRunAddress, "адрес и порт запуска сервиса")
	flagDatabaseURI := flag.String("d", defaultDatabaseURI, "адрес подключения к базе данных")
	flagAccrualSystemAddress := flag.String("r", defaultAccrualSystemAddress, "адрес системы расчёта начислений")

	// Парсим флаги
	flag.Parse()

	// Если переменные окружения заданы, они переопределяют значения по умолчанию
	if envRunAddress != "" {
		*flagRunAddress = envRunAddress
	}
	if envDatabaseURI != "" {
		*flagDatabaseURI = envDatabaseURI
	}
	if envAccrualSystemAddress != "" {
		*flagAccrualSystemAddress = envAccrualSystemAddress
	}

	// Возвращаем конфигурацию
	return global.ConfigFlag{
		RunAddress:           *flagRunAddress,
		DatabaseURI:          *flagDatabaseURI,
		AccrualSystemAddress: *flagAccrualSystemAddress,
	}
}
