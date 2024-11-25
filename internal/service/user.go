package service

import (
	"encoding/json"
	"fmt"
	"github.com/LI-SeNyA-vE/Gofermat/internal/global"
	"io"
	"net/http"
	"strconv"
	"time"
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
	_, err = strconv.Atoi(orderNumber)
	if err != nil {
		return 422, fmt.Errorf("введены не только цифры %s", err)
	}
	for i := len(orderNumber) - 1; i >= 0; i-- {
		r := rune(orderNumber[i])

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

	go sendOrderToAPI(orderUser.NumberOrder)

	return statusCode, err
}

type OrderResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func testSendOrderToAPI(numberOrder string) {
	http.HandleFunc("/api/orders/", func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем номер заказа из URL
		orderNumber := r.URL.Path[len("/api/orders/"):]

		// Формируем фиксированный ответ
		response := OrderResponse{
			Order:   orderNumber,
			Status:  "PROCESSED",
			Accrual: 500,
		}

		// Устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "application/json")

		// Возвращаем JSON-ответ
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		}
	})
}
func sendOrderToAPI(numberOrder string) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/orders/%s", global.Config.Flags.AccrualSystemAddress, numberOrder)

	for {
		// Отправляем GET-запрос
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Ошибка запроса: %v. Повтор через 10 секунд...\n", err)
			time.Sleep(10 * time.Second)
			continue
		}

		defer resp.Body.Close()

		// Проверяем статус ответа
		if resp.StatusCode != http.StatusOK {
			continue
		}

		// Читаем тело ответа
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}

		var orderFromAccrualSystem global.OrderWithdrawalsUserJSON

		err = json.Unmarshal(body, &orderFromAccrualSystem)
		if err != nil {
			return
		}

		err = UpdateOrder(orderFromAccrualSystem)
		if err != nil {
			return
		}
		break
	}
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
