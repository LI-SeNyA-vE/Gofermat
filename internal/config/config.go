package config

import (
	"flag"
	"github.com/LI-SeNyA-vE/Gofermat/internal/logger"
	"github.com/LI-SeNyA-vE/Gofermat/internal/model"
	"os"
)

// WrapperConfig обёртка для model.Configs с методами
type WrapperConfig struct {
	Configs *model.Configs
}

func Start() (*model.Configs, error) {
	err := logger.Initialize("info")
	if err != nil {
		return nil, err
	}
	conf := newConfig()
	conf.loadConfig()
	return conf.Configs, nil
}

func newConfig() *WrapperConfig {
	return &WrapperConfig{}
}

func (conf *WrapperConfig) loadConfig() {
	// Устанавливаем значения по умолчанию
	defaultRunAddress := "localhost:8080"
	defaultDatabaseURI := "postgres://senya:1q2w3e4r5t@localhost:5433/gofermat"
	defaultAccrualSystemAddress := "http://localhost:8081"
	defaultSecretKeyForJWT := "SecretKeyForJWT"
	defaultSecretKeyForPassword := "SecretKeyForPassword"

	// Читаем переменные окружения
	envRunAddress := os.Getenv("RUN_ADDRESS")
	envDatabaseURI := os.Getenv("DATABASE_URI")
	envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	envSecretKeyForJWT := os.Getenv("SECRET_KEY_FOR_JWT")
	envSecretKeyForPassword := os.Getenv("SECRET_KEY_FOR_PASSWORD")

	// Определяем флаги
	flagRunAddress := flag.String("a", defaultRunAddress, "адрес и порт запуска сервиса")
	flagDatabaseURI := flag.String("d", defaultDatabaseURI, "адрес подключения к базе данных")
	flagAccrualSystemAddress := flag.String("r", defaultAccrualSystemAddress, "адрес системы расчёта начислений")
	flagSecretKeyForJWT := flag.String("j", defaultSecretKeyForJWT, "секретный ключ для создания JWT")
	flagSecretKeyForPassword := flag.String("w", defaultSecretKeyForPassword, "секретный ключ для HASH пароля")

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

	if envSecretKeyForJWT != "" {
		*flagSecretKeyForJWT = envSecretKeyForJWT
	}

	if envSecretKeyForPassword != "" {
		*flagSecretKeyForPassword = envSecretKeyForPassword
	}

	// Возвращаем конфигурацию

	conf.Configs.ConfigFlag = model.ConfigFlag{
		RunAddress:           *flagRunAddress,
		DatabaseURI:          *flagDatabaseURI,
		AccrualSystemAddress: *flagAccrualSystemAddress,
		SecretKeyForJWT:      *flagSecretKeyForJWT,
		SecretKeyForPassword: *flagSecretKeyForPassword,
	}
}
