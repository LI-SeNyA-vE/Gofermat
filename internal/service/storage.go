package service

import (
	"fmt"
	"github.com/LI-SeNyA-vE/Gofermat/internal/config"
	"github.com/LI-SeNyA-vE/Gofermat/internal/global"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	log = global.Logger
)

type Storage struct {
	Config *config.Config
	User   *User
}

func NewStorage() *Storage {
	return &Storage{
		Config: config.New(),
		User:   NewUser(),
	}
}

// connectToDataBase Служит для подключения к базе
// Если нет передаваемого значения, то по умолчанию идёт подключения к postgres
func (storage *Storage) connectToDataBase(dataBase ...string) (*gorm.DB, error) {
	var nameDataBase string
	if len(dataBase) == 0 {
		nameDataBase = storage.Config.DataBase.NameDataBase
	} else {
		nameDataBase = dataBase[0]
	}
	connectBD := fmt.Sprintf("host=%s dbname=%s port=%s User=%s password=%s sslmode=disabl",
		storage.Config.DataBase.Host,
		nameDataBase,
		storage.Config.DataBase.Port,
		storage.Config.DataBase.User,
		storage.Config.DataBase.Password,
	)

	db, err := gorm.Open(postgres.Open(connectBD), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Infof("Ошибка подключения к базе данных%s: %v", nameDataBase, err)
		return nil, err
	}
	return db, nil
}

func (storage *Storage) createDataBaseIfNotExists() error {
	gormDB, err := storage.connectToDataBase()
	if err != nil {
		return err
	}

	defer func() {
		sqlDB, _ := gormDB.DB()
		_ = sqlDB.Close()
	}()

	exists := gormDB.Exec("SELECT 1 FROM pg_database WHERE datname = ?", storage.Config.DataBase.NameDataBase).Error
	if exists != nil {
		err = gormDB.Exec("CREATE DATABASE ?" + storage.Config.DataBase.NameDataBase).Error
		if err != nil {
			return fmt.Errorf("не удалось создать базу данных: %w", err)
		}
		log.Infof("База данных %s успешно создана", storage.Config.DataBase.NameDataBase)
	} else {
		log.Infof("База данных %s уже существует", storage.Config.DataBase.NameDataBase)
	}
	return nil
}

func (storage *Storage) autoMigrateTables(db *gorm.DB) error {
	models := []interface{}{
		storage.User.UserCred,
	}
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("не удалось создать таблицу: %w", err)
		}
	}
	return nil
}
func (storage *Storage) Start() error {
	// Создаём базу данных, если она отсутствует
	if err := storage.createDataBaseIfNotExists(); err != nil {
		log.Fatalf("Ошибка при создании базы данных: %v", err)
		return err
	}

	// Подключаемся к базе данных
	gormDB, err := storage.connectToDataBase(storage.Config.DataBase.NameDataBase)
	defer func() {
		sqlDB, _ := gormDB.DB()
		_ = sqlDB.Close()
	}()

	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
		return err
	}

	// Выполняем миграцию для создания таблиц
	err = storage.autoMigrateTables(gormDB)
	if err != nil {
		log.Fatalf("Ошибка при миграции таблиц: %v", err)
		return err
	}

	log.Infof("Подключение к базе данных установлено и таблицы проверены/созданы.")
	return nil
}
