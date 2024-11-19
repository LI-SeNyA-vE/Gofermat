package service

import (
	"encoding/json"
	"fmt"
	"github.com/LI-SeNyA-vE/Gofermat/internal/global"
	"net/http"
	"strconv"
)

func CreateUser(userCred global.UserCred) (token string, statusCode int, err error) {
	token, statusCode, err = RegistrationUser(&userCred)
	return token, statusCode, err
}

func UserAuthentication(userCred global.UserCred) (token string, statusCode int, err error) {
	token, statusCode, err = AuthenticationUser(&userCred)
	return token, statusCode, err
}

func UserUploadingNumberOrder(jwt string, orderNumber int) (statusCode int, err error) {
	//реализация алгоритма Луна
	// Перебираем число в обратном порядке
	var sum int
	var double = false
	orderNumberString := strconv.Itoa(orderNumber)
	for i := len(orderNumberString) - 1; i >= 0; i-- {
		r := rune(orderNumberString[i])

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
	//конец реализации

	payload, _ := VerificationToken(jwt)
	userLogin := payload["sub"].(string)
	userId, statusCode, err := SearchForUserById(userLogin)
	if err != nil {
		return statusCode, fmt.Errorf("ошибка на поиске пользователя %v", err)
	}
	var orderUser = global.OrderUser{UserId: userId, NumberOrder: orderNumber}
	statusCode, err = UploadingNumberOrder(orderUser)

	return statusCode, err
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
		return ordersJSON, http.StatusInternalServerError, fmt.Errorf("на этапе маршла данных произошла ошибка %v", err)
	}

	return ordersJSON, http.StatusOK, nil
}

func UserListUserBalance(jwt string) (balanceJSON []byte, statusCode int, err error) {

	payload, _ := VerificationToken(jwt)
	userLogin := payload["sub"].(string)
	balance, statusCode, err := GetUserBalance(userLogin)
	if err != nil {
		return balanceJSON, statusCode, fmt.Errorf("на этапе получения данных из базы произошла ошибка %v", err)
	}

	balanceJSON, err = json.Marshal(balance)
	if err != nil {
		return nil, 500, err
	}

	return nil, 500, err
}
