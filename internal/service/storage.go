package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/LI-SeNyA-vE/Gofermat/internal/global"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
)

// Hash создаёт хеш для паролей
func hash(src string) []byte {
	h := hmac.New(sha256.New, []byte(global.Config.HashKey))
	h.Write([]byte(src))
	dst := h.Sum(nil)
	return dst
}

func SearchForUserById(userLogin string) (userId uint, statusCode int, err error) {
	var user global.UserCredInDataBase
	gormDB, err := ConnectToDataBase()
	if err != nil {
		return 0, http.StatusInternalServerError, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	err = gormDB.Where("login = ?", userLogin).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, http.StatusUnauthorized, fmt.Errorf("пользователь с login: %s не найден: %w", userLogin, gorm.ErrRecordNotFound)
		}
		return 0, http.StatusInternalServerError, err
	}
	return user.ID, http.StatusOK, nil
}

// ConnectToDataBase Служит для подключения к базе
// Если нет передаваемого значения, то по умолчанию идёт подключения к postgres
func ConnectToDataBase(dataBase ...string) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(global.Config.Flags.DatabaseURI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		global.Logger.Infof("Ошибка подключения к базе данных: %v", err)
		return nil, err
	}
	return db, nil
}

// CreateDataBaseIfNotExists создаёт базу если она отсутствует
func CreateDataBaseIfNotExists() error {
	gormDB, err := ConnectToDataBase("postgres")
	if err != nil {
		global.Logger.Infof("ошибка при подключение к базе postgres %v", err)
		return err
	}

	defer func() {
		sqlDB, _ := gormDB.DB()
		_ = sqlDB.Close()
	}()
	var exists bool
	err = gormDB.Raw("SELECT 1 FROM pg_database WHERE datname = ?", global.Config.DataBase.NameDataBase).Scan(&exists).Error
	if err != nil {
		global.Logger.Infof("ошибка при проверке существования базы данных: %v", err)
		return err
	}
	if !exists {
		err = gormDB.Exec("CREATE DATABASE " + global.Config.DataBase.NameDataBase).Error
		if err != nil {
			return fmt.Errorf("не удалось создать базу данных: %w", err)
		}
		global.Logger.Infof("База данных %s успешно создана", global.Config.DataBase.NameDataBase)
	} else {
		global.Logger.Infof("База данных %s уже существует", global.Config.DataBase.NameDataBase)
	}
	return nil
}

// StartStorage использоваться при запуске программы и передаёт обратно переменную для подключения
func StartStorage() (*gorm.DB, error) {
	// Создаём базу данных, если она отсутствует
	//if err := CreateDataBaseIfNotExists(); err != nil {
	//	log.Fatalf("Ошибка при создании базы данных: %v", err)
	//	return nil, err
	//}

	// Подключаемся к базе данных
	gormDB, err := ConnectToDataBase()

	if err != nil {
		global.Logger.Infof("Ошибка подключения к базе данных: %v", err)
		return nil, err
	}

	err = gormDB.AutoMigrate(global.UserCredInDataBase{}, global.OrderUser{}, global.BalanceUser{})
	if err != nil {
		global.Logger.Infof("ошибка %v при создании таблиц", err) //ошибка на миграции базы
	}

	//err = gormDB.Unscoped().Delete(&global.UserCredInDataBase{}, 2).Error //Это удаление строк в таблице

	global.Logger.Infof("Подключение к базе данных установлено")
	return gormDB, nil
}

// searchUserIBase делает поиск в таблице пользователя с переданным логином
func searchUserIBase(gormDB *gorm.DB, login string) (user *global.UserCredInDataBase, err error) {
	err = gormDB.Where("login = ?", login).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("пользователь с login: %s не найден: %w", login, gorm.ErrRecordNotFound)
		}
		return nil, err
	}
	return user, nil
}

// searchOrderIBase делает поиск в таблице пользователя с переданным логином
// если ошибка order!=nil и errors.Is(err, gorm.ErrRecordNotFound), значит заказ с таким номером ещё не загружали
// если order!=nil значит такой заказ уже есть в базе
func searchOrderIBase(gormDB *gorm.DB, numberOrder string) (order *global.OrderUser, err error) {
	err = gormDB.Where("number_order = ?", numberOrder).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("такого заказа ещё нет: %w", err)
		}
		return nil, fmt.Errorf("ошибка при поиске заказа в таблице: %w", err)
	}
	return order, fmt.Errorf("такой заказ уже загружен")
}

// Регистрирует пользователя.
// Проверяет, существует пользователь или нет
func RegistrationUser(userCred *global.UserCred) (jwt string, statusCode int, err error) {
	var token string
	gormDB, err := ConnectToDataBase()
	if err != nil {
		return token, http.StatusInternalServerError, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	var userCredInDataBase = global.UserCredInDataBase{
		Login:    userCred.Login,
		Password: hash(userCred.Password),
	}

	userIBase, err := searchUserIBase(gormDB, userCredInDataBase.Login)
	//если пользователь не найден
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return token, http.StatusInternalServerError, fmt.Errorf("ошибка %v при запросе к бизе данных для получения данных пользователя", err) //Если ошибка произошла на этапе поиска/запроса в базу
		}
	}

	if userIBase != nil {
		return token, http.StatusConflict, fmt.Errorf("пользователь с таким логином: %v уже существует", userCredInDataBase.Login)
	}

	gormDB.Create(&userCredInDataBase)
	gormDB.Create(&global.BalanceUser{
		UserId:    userCredInDataBase.ID,
		Current:   0,
		Withdrawn: 0,
	})

	token, err = CreateNewToken(userCred.Login)
	if err != nil {
		return token, http.StatusInternalServerError, err
	}

	return token, http.StatusOK, nil
}

// AuthenticationUser авторизовывает пользователя (проверяет, существует он в талице или нет) в ответ передаёт токен
func AuthenticationUser(userCred *global.UserCred) (jwt string, statusCode int, err error) {
	var token string
	gormDB, err := ConnectToDataBase()
	if err != nil {
		return token, http.StatusInternalServerError, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	var userCredInDataBase = global.UserCredInDataBase{
		Login:    userCred.Login,
		Password: hash(userCred.Password),
	}

	userIBase, err := searchUserIBase(gormDB, userCredInDataBase.Login)
	//если пользователь не найден
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return token, http.StatusInternalServerError, fmt.Errorf("ошибка %v при запросе к бизе данных для получения данных пользователя", err) //Если ошибка произошла на этапе поиска/запроса в базу
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return token, http.StatusUnauthorized, fmt.Errorf("пользователь с таким логином: %v, не найден", userCredInDataBase.Login) //пользователь не найден
	}

	if !bytes.Equal(userIBase.Password, userCredInDataBase.Password) {
		return token, http.StatusUnauthorized, fmt.Errorf("введён неверный пароль") //пользователь ввёл неверные данные
	}

	token, err = CreateNewToken(userCred.Login)
	if err != nil {
		return token, http.StatusInternalServerError, fmt.Errorf("ошибка при создании токена")
	}
	return token, http.StatusOK, nil //пользователь ввёл правильный пароль
}

// UploadingNumberOrder загружает в базу заказ пользователя если его нет, либо информацию об этом заказе (за кем он уже зарегистрирован)
func UploadingNumberOrder(orderUser global.OrderUser) (statusCode int, err error) {
	gormDB, err := ConnectToDataBase()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	orderUserFromDB, err := searchOrderIBase(gormDB, orderUser.NumberOrder)

	if orderUserFromDB != nil {
		if orderUserFromDB.UserId == orderUser.UserId {
			return http.StatusOK, fmt.Errorf("номер заказа уже был загружен этим пользователем")
		} else {
			return http.StatusConflict, fmt.Errorf("номер заказа уже был загружен другим пользователем")
		}
	}
	gormDB.Create(&orderUser)

	return http.StatusAccepted, nil
}

func GetListUserOrders(userId uint) (userOrder []global.OrderUser, statusCode int, err error) {
	gormDB, err := ConnectToDataBase()
	if err != nil {
		return []global.OrderUser{}, http.StatusInternalServerError, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	gormDB.Where("user_id = ?", userId).Find(&userOrder)

	if len(userOrder) == 0 {
		return nil, http.StatusNoContent, fmt.Errorf("нет данных для ответа: %v", err)
	}

	return userOrder, 500, nil
}

func GetUserBalance(userId uint) (userBalance global.BalanceUser, statusCode int, err error) {
	gormDB, err := ConnectToDataBase()
	if err != nil {
		return global.BalanceUser{}, http.StatusInternalServerError, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	err = gormDB.Where("user_id = ?", userId).First(&userBalance).Error
	if err != nil {
		return global.BalanceUser{}, http.StatusInternalServerError, fmt.Errorf("ошибка при получение баланса пользователя из табилцы %v", err)
	}
	return userBalance, 200, err
}

func DebtSumFromBalanceAndCreateOrders(orderForPoints global.OrderForPoints, userId uint) (statusCode int, err error) {
	var orderUser = global.OrderUser{
		UserId:      userId,
		NumberOrder: orderForPoints.NumberOrder,
		Status:      "",
		Accrual:     0,
		Sum:         orderForPoints.Sum,
		Model:       global.Model{},
	}

	gormDB, err := ConnectToDataBase()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	balance, statusCode, err := GetUserBalance(userId)
	if err != nil {
		return statusCode, fmt.Errorf("на этапе получения данных из базы произошла ошибка %v", err)
	}

	balance.Current -= orderForPoints.Sum
	balance.Withdrawn += orderForPoints.Sum

	if err = gormDB.Save(&balance).Error; err != nil {
		return 500, fmt.Errorf("ошибка при обновлении баланса: %w", err)
	}

	statusCode, err = UploadingNumberOrder(orderUser)
	if err != nil {
		return statusCode, err
	}

	return statusCode, nil
}

func GetOrdersWithdrawal(userId uint) (usersOrder []global.OrderUser, statusCode int, err error) {
	gormDB, err := ConnectToDataBase()
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	err = gormDB.Where("user_id = ? AND sum IS NOT NULL AND sum != 0", userId).Find(&usersOrder).Error
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("ошибка получения данных из таблицы: %v", err)
	}
	if len(usersOrder) == 0 {
		return nil, 204, nil
	}

	return usersOrder, 200, nil
}

func UpdateOrder(orderFromAccrualSystem global.OrderWithdrawalsUserJSON) error {
	gormDB, err := ConnectToDataBase()
	if err != nil {
		return err
	}

	if err := gormDB.Model(global.OrderUser{}).
		Where("number_order = ?", orderFromAccrualSystem.NumberOrder).
		Updates(map[string]interface{}{
			"status":  orderFromAccrualSystem.Status,
			"accrual": orderFromAccrualSystem.Accrual,
		}).Error; err != nil {
		return fmt.Errorf("ошибка при обновлении нескольких полей: %w", err)
	}
	return nil

}
