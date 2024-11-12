package config

import (
	"github.com/LI-SeNyA-vE/Gofermat/internal/global"
	"gopkg.in/yaml.v3"
	"os"
)

var logger = global.Logger

type Config struct {
	DataBase *DataBase
}

type DataBase struct {
	Host         string `yaml:"POSTGRES_HOST"`
	User         string `yaml:"POSTGRES_USER"`
	Password     string `yaml:"POSTGRES_PASSWORD"`
	Port         string `yaml:"PORT"`
	NameDataBase string `yaml:"POSTGRES_DB"`
}

func New() *Config {
	return &Config{}
}

func (config *Config) LoadConfigFromFile(fileName string) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		logger.Errorf("ошибка чтения конфигурационного файла %s\bОшибка: %v", fileName, err)
	}

	err = yaml.Unmarshal(data, &config.DataBase)
	if err != nil {
		logger.Errorf("ошибка парсинга Yaml файла в config\b%v", err)
	}
}
