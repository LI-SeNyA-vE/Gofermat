package global

import (
	"gorm.io/gorm"
	"time"
)

// Для пакета Connfig
type Configs struct {
	DataBase  DataBaseConfig `yaml:"POSTGRES_CONFIG"`
	SecretKey string         `yaml:"SECRET_KEY"`
	HashKey   string         `yaml:"HASH_KEY"`
}

type DataBaseConfig struct {
	Host         string `yaml:"POSTGRES_HOST"`
	User         string `yaml:"POSTGRES_USER"`
	Password     string `yaml:"POSTGRES_PASSWORD"`
	Port         string `yaml:"PORT"`
	NameDataBase string `yaml:"POSTGRES_DB"`
}

// Для пакета storage
type Model struct {
	ID        uint           `gorm:"primary_key" json:"-"`
	CreatedAt time.Time      `gorm:"uploaded_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserCredInDataBase struct {
	Login    string `gorm:"unique;index"`
	Password []byte `gorm:"BINARY(32)" json:"-"`
	Model
	Orders []OrderUser `gorm:"foreignKey:UserId;references:ID"`
}

type OrderUser struct {
	UserId      uint    `gorm:"index" json:"-"`
	NumberOrder int     `gorm:"unique" json:"order"`
	Status      string  `gorm:"DEFAULT:NEW"`
	Accrual     float32 `json:"Accrual,omitempty"`
	Model
}

type BalanceUser struct {
	current   float32
	withdrawn float32
	Model
}

// Для пакета user
type UserCred struct {
	Id       int64  `gorm:"primary_key" json:"id"`
	Login    string `gorm:"unique" json:"login"`
	Password string `json:"password"`
}

// Для пакета jwt

// Для пакета

// Для пакета

// Для пакета

// Для пакета
