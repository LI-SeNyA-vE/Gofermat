package global

import (
	"gorm.io/gorm"
	"time"
)

// Для пакета Connfig
type Configs struct {
	DataBase  ConfigDataBase `yaml:"POSTGRES_CONFIG"`
	Flags     ConfigFlag
	SecretKey string `yaml:"SECRET_KEY"`
	HashKey   string `yaml:"HASH_KEY"`
}

type ConfigFlag struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
}

type ConfigDataBase struct {
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
	Orders  []OrderUser `gorm:"foreignKey:UserId;references:ID;constraint:OnDelete:CASCADE"`
	Balance BalanceUser `gorm:"foreignKey:UserId;references:ID;constraint:OnDelete:CASCADE"`
}

type OrderUser struct {
	UserId      uint    `gorm:"index" json:"-"`
	NumberOrder string  `gorm:"unique" json:"order"`
	Status      string  `gorm:"DEFAULT:NEW"`
	Accrual     float32 `json:"accrual,omitempty"`
	Sum         float32 `json:"sum,omitempty"`
	Model
}

type BalanceUser struct {
	ID        uint `gorm:"primary_key" json:"-"`
	UserId    uint `gorm:"unique index" json:"-"`
	Current   float32
	Withdrawn float32
	CreatedAt time.Time      `gorm:"uploaded_at" json:"-"`
	UpdatedAt time.Time      `gorm:"processed_at" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderWithdrawalsUserJSON struct {
	UserId      uint           `json:"-"`
	NumberOrder string         `json:"order"`
	Status      string         `json:"-"`
	Accrual     float32        `json:"-"`
	Sum         float32        `json:"sum"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"processed_at"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}

// Для пакета user
type UserCred struct {
	Id       uint   `gorm:"primary_key" json:"id"`
	Login    string `gorm:"unique" json:"login"`
	Password string `json:"password"`
}

// Для пакета jwt

// Для пакета handlers

type OrderForPoints struct {
	NumberOrder string  `json:"order"`
	Sum         float32 `json:"sum"`
}

// Для пакета

// Для пакета

// Для пакета
