package service

import (
	"encoding/json"
	"fmt"
	"github.com/LI-SeNyA-vE/Gofermat/internal/global"
	"io"
	"net/http"
	"strconv"
	"unicode"
)

func CreateUser(userCred global.UserCred) (token string, statusCode int, err error) {
	token, statusCode, err = RegistrationUser(&userCred)
	return token, statusCode, err
}

func UserAuthentication(userCred global.UserCred) (token string, statusCode int, err error) {
	token, statusCode, err = AuthenticationUser(&userCred)
	return token, statusCode, err
}

func LunaAlgorithm(orderNumber string) (statusCode int, err error) {
	var sum int
	var double = false
	if len(orderNumber) == 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf("номер заказа не прошёл проверку по алгоритму \"Луна\" %w", err)
	}
	for i := len(orderNumber) - 1; i >= 0; i-- {
		r := rune(orderNumber[i])
		if !unicode.IsDigit(r) {
			return http.StatusUnprocessableEntity, fmt.Errorf("введены не только цифры %s", err)
		}
		digit, _ := strconv.Atoi(string(r))

		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		double = !double
	}

	if sum%10 != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf("номер заказа не прошёл проверку по алгоритму \"Луна\" %w", err)
	}
	return http.StatusOK, nil
}

func UserUploadingNumberOrder(jwt string, orderNumber string) (statusCode int, err error) {
	statusCode, err = LunaAlgorithm(orderNumber)
	if err != nil {
		return statusCode, err
	}
	payload, _ := VerificationToken(jwt)
	userLogin := payload["sub"].(string)
	userId, statusCode, err := SearchForUserById(userLogin)
	if err != nil {
		return statusCode, fmt.Errorf("ошибка на поиске пользователя %v", err)
	}
	var orderUser = global.OrderUser{UserId: userId, NumberOrder: orderNumber}
	statusCode, err = UploadingNumberOrder(orderUser)

	if statusCode != http.StatusOK && statusCode != http.StatusConflict {
		go sendOrderToAPI(orderUser.NumberOrder)
	}

	return statusCode, err
}

type OrderResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func sendOrderToAPI(numberOrder string) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/orders/%s", global.Config.Flags.AccrualSystemAddress, numberOrder)
	global.Logger.Infof("url = %s", url)

	// Отправляем GET-запрос
	resp, err := http.Get(url)
	if err != nil {
		global.Logger.Infof("Ошибка запроса: %v.  К системе лояльности...\n", err)
		return
	}

	defer resp.Body.Close()

	// Проверяем статус ответа
	switch resp.StatusCode {
	case http.StatusOK:
		global.Logger.Infof("resp.StatusCode = http.StatusOK в системе лояльности")
	case http.StatusNoContent:
		global.Logger.Infof("resp.StatusCode = http.StatusNoContent в системе лояльности")
	case http.StatusInternalServerError:
		global.Logger.Infof("resp.StatusCode = http.StatusInternalServerError в системе лояльности")
	case http.StatusTooManyRequests:
		timeSlip := resp.Header.Get("Retry-After")
		intTimeSlip, _ := strconv.Atoi(timeSlip)
		global.Logger.Infof("resp.StatusCode = http.StatusInternalServerError в системе лояльности\bпревышено колличество запросов, подождите %v", intTimeSlip)
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		global.Logger.Infof("Ошибка при чтение ответа в системе лояльности, %v", err)
		return
	}

	var orderFromAccrualSystem global.OrderWithdrawalsUserJSON

	test := string(body)
	global.Logger.Infof("test = %s", test)

	err = json.Unmarshal(body, &orderFromAccrualSystem)
	if err != nil {
		global.Logger.Infof("Ошибка при разбре body в orderFromAccrualSystem global.OrderWithdrawalsUserJSON в системе лояльности, %v", err)
		return
	}

	err = UpdateOrder(orderFromAccrualSystem)
	if err != nil {
		global.Logger.Infof("Ошибка при записи в базу в системе лояльности, %v", err)
		return
	}

	global.Logger.Infof("обратились к системе лояльности %v", orderFromAccrualSystem)
}

func UserListUserOrders(jwt string) (ordersJSON []byte, statusCode int, err error) {
	var orders []global.OrderUser
	payload, _ := VerificationToken(jwt)
	userLogin := payload["sub"].(string)
	userId, statusCode, err := SearchForUserById(userLogin)
	if err != nil {
		return nil, statusCode, fmt.Errorf("ошибка на поиске пользователя %v", err)
	}

	orders, statusCode, err = GetListUserOrders(userId)
	if err != nil {
		return ordersJSON, statusCode, fmt.Errorf("на этапе получения данных из базы произошла ошибка %v", err)
	}
	ordersJSON, err = json.Marshal(orders)
	if err != nil {
		return ordersJSON, 500, fmt.Errorf("на этапе маршла данных произошла ошибка %v", err)
	}

	return ordersJSON, http.StatusOK, nil
}

func UserListUserBalance(jwt string) (balanceJSON []byte, statusCode int, err error) {

	payload, _ := VerificationToken(jwt)
	userLogin := payload["sub"].(string)
	userId, statusCode, err := SearchForUserById(userLogin)
	if err != nil {
		return nil, statusCode, fmt.Errorf("ошибка на поиске пользователя %v", err)
	}
	balance, statusCode, err := GetUserBalance(userId)
	if err != nil {
		return balanceJSON, statusCode, fmt.Errorf("на этапе получения данных из базы произошла ошибка %v", err)
	}

	balanceJSON, err = json.Marshal(balance)
	if err != nil {
		return nil, 500, err
	}

	return balanceJSON, statusCode, nil
}

func UserNewOrderForPoints(orderForPoints global.OrderForPoints, jwt string) (statusCode int, err error) {
	statusCode, err = LunaAlgorithm(orderForPoints.NumberOrder)
	if err != nil {
		return statusCode, err
	}

	payload, _ := VerificationToken(jwt)
	userLogin := payload["sub"].(string)
	userId, statusCode, err := SearchForUserById(userLogin)
	if err != nil {
		return statusCode, fmt.Errorf("ошибка на поиске пользователя %v", err)
	}
	balance, statusCode, err := GetUserBalance(userId)
	if err != nil {
		return statusCode, fmt.Errorf("на этапе получения данных из базы произошла ошибка %v", err)
	}

	if orderForPoints.Sum <= balance.Current {
		statusCode, err = DebtSumFromBalanceAndCreateOrders(orderForPoints, userId)
		if err != nil {
			return statusCode, err
		}
	} else {
		return 402, fmt.Errorf("недостаточко средств")
	}
	return statusCode, nil
}

func OrdersPaidPoints(jwt string) (userOrdersJSON []global.OrderWithdrawalsUserJSON, statusCode int, err error) {
	payload, _ := VerificationToken(jwt)
	userLogin := payload["sub"].(string)

	userId, statusCode, err := SearchForUserById(userLogin)
	if err != nil {
		return nil, statusCode, fmt.Errorf("ошибка на поиске пользователя %v", err)
	}

	userOrders, statusCode, err := GetOrdersWithdrawal(userId)
	if err != nil {
		return nil, 0, err
	}

	for _, userOrder := range userOrders {
		userOrdersJSON = append(userOrdersJSON, global.OrderWithdrawalsUserJSON{
			UserId:      userOrder.UserId,
			NumberOrder: userOrder.NumberOrder,
			Status:      userOrder.Status,
			Accrual:     userOrder.Accrual,
			Sum:         userOrder.Sum,
			CreatedAt:   userOrder.CreatedAt,
			UpdatedAt:   userOrder.UpdatedAt,
			DeletedAt:   userOrder.DeletedAt,
		})
	}

	return userOrdersJSON, statusCode, nil
}
