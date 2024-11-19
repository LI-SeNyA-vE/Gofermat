package service

import "gorm.io/gorm"

func Start() (*gorm.DB, error) {
	//запуск базы данных
	gormDB, err := StartStorage()
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}
