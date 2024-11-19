package config

import (
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
