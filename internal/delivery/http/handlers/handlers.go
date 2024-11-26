package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/LI-SeNyA-vE/Gofermat/internal/global"
	"github.com/LI-SeNyA-vE/Gofermat/internal/service"
	"net/http"
	"strings"
)

func UserRegistration(writer http.ResponseWriter, request *http.Request) {
	var (
		buf             bytes.Buffer
		userCredentials global.UserCred
	)

	if request.Header.Get("Content-Type") != "application/json" {
		global.Logger.Info("неверный формат запроса")
		http.Error(writer, fmt.Sprint("неверный формат запроса"), http.StatusBadRequest)
		return
	}

	_, err := buf.ReadFrom(request.Body) //Читает данные из тела запроса
	if err != nil {
		global.Logger.Info("ошибка при чтении данных из Body")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &userCredentials) // Разбирает данные из массива byte в структуру
	if err != nil || userCredentials.Password == "" || userCredentials.Login == "" {
		global.Logger.Info("неверный формат данных переданный пользователем")
		http.Error(writer, fmt.Sprintf("неверный формат данных переданный пользователем"), http.StatusBadRequest)
		return
	}

	token, statusCode, err := service.CreateUser(userCredentials) // Функция на создание пользователя
	if err != nil {
		global.Logger.Info(err)
		http.Error(writer, err.Error(), statusCode)
		return
	}
	writer.Header().Set("Authorization", "Bearer "+token)
	writer.WriteHeader(http.StatusOK)
}

func UserAuthentication(writer http.ResponseWriter, request *http.Request) {
	var (
		buf             bytes.Buffer
		userCredentials global.UserCred
	)

	if request.Header.Get("Content-Type") != "application/json" {
		global.Logger.Info("неверный формат запроса")
		http.Error(writer, fmt.Sprint("неверный формат запроса"), http.StatusBadRequest)
		return
	}

	_, err := buf.ReadFrom(request.Body) //Читает данные из тела запроса
	if err != nil {
		global.Logger.Info("ошибка при чтении данных из Body")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &userCredentials) // Разбирает данные из массива byte в структуру
	if err != nil || userCredentials.Password == "" || userCredentials.Login == "" {
		global.Logger.Info("неверный формат данных переданный пользователем")
		http.Error(writer, fmt.Sprintf("неверный формат данных переданный пользователем"), http.StatusBadRequest)
		return
	}

	token, statusCode, err := service.UserAuthentication(userCredentials) // Функция на Авторизацию пользователя
	if err != nil {
		global.Logger.Info(err)
		http.Error(writer, err.Error(), statusCode)
		return
	}
	writer.Header().Set("Authorization", "Bearer "+token)
	writer.WriteHeader(http.StatusOK)
}

func AddOrder(writer http.ResponseWriter, request *http.Request) {
	var (
		buf         bytes.Buffer
		numberOrder string
	)

	if request.Header.Get("Content-Type") != "text/plain" {
		global.Logger.Info("неверный формат запроса")
		http.Error(writer, fmt.Sprint("неверный формат запроса"), http.StatusBadRequest)
		return
	}

	_, err := buf.ReadFrom(request.Body) //Читает данные из тела запроса
	if err != nil {
		global.Logger.Info("ошибка при чтении данных из Body")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	numberOrder = buf.String()
	global.Logger.Infof("заказ № %s", numberOrder)

	authHeader := request.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	global.Logger.Infof("токен передающийся в функцию = %s", token)
	statusCode, err := service.UserUploadingNumberOrder(token, numberOrder)
	if err != nil {
		global.Logger.Infof("ошибка на этапе загрузки пользователем заказа %s", err)
		http.Error(writer, err.Error(), statusCode)
		return
	}

	writer.WriteHeader(statusCode)
}

func ExpenditurePointsOnNewOrder(writer http.ResponseWriter, request *http.Request) {
	var (
		buf             bytes.Buffer
		userCredentials global.OrderForPoints
	)

	if request.Header.Get("Content-Type") != "application/json" {
		global.Logger.Info("неверный формат запроса")
		http.Error(writer, fmt.Sprint("неверный формат запроса"), http.StatusBadRequest)
		return
	}

	_, err := buf.ReadFrom(request.Body) //Читает данные из тела запроса
	if err != nil {
		global.Logger.Info("ошибка при чтении данных из Body")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &userCredentials) // Разбирает данные из массива byte в структуру
	if err != nil {
		global.Logger.Info("неверный формат данных")
		http.Error(writer, fmt.Sprintf("неверный формат данных"), http.StatusBadRequest)
		return
	}

	authHeader := request.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	global.Logger.Infof("токен передающийся в функцию = %s", token)
	statusCode, err := service.UserNewOrderForPoints(userCredentials, token)
	if err != nil {
		global.Logger.Info("ошибка при выполнение функции UserNewOrderForPoints")
		http.Error(writer, err.Error(), statusCode)
		return
	}

	writer.WriteHeader(statusCode)
	request.Header.Get("Заказ успешно зарегистрирован и списаны баллы")

}

func ListUserOrders(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("Content-Length") != "0" {
		http.Error(writer, "Content-Length != 0", http.StatusBadRequest)
		return
	}

	authHeader := request.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	global.Logger.Infof("токен передающийся в функцию = %s", token)
	ordersJSON, statusCode, err := service.UserListUserOrders(token)
	if err != nil {
		global.Logger.Infof("ошибка на этапе выгрузке заказов пользователя %s", err)
		http.Error(writer, err.Error(), statusCode)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(ordersJSON)
	writer.WriteHeader(http.StatusOK)
}

func ListUserBalance(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("Content-Length") != "0" {
		http.Error(writer, "Content-Length != 0", http.StatusBadRequest)
		return
	}

	authHeader := request.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	global.Logger.Infof("токен передающийся в функцию = %s", token)
	balanceJSON, statusCode, err := service.UserListUserBalance(token)
	if err != nil {
		global.Logger.Infof("ошибка на этапе выгрузке баланса пользователя %s", err)
		http.Error(writer, err.Error(), statusCode)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(balanceJSON)
	writer.WriteHeader(http.StatusOK)
}

func InfoAboutUsagePoints(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("Content-Length") != "0" {
		http.Error(writer, "Content-Length != 0", http.StatusBadRequest)
		return
	}

	authHeader := request.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	global.Logger.Infof("токен передающийся в функцию = %s", token)
	usersOrder, statusCode, err := service.OrdersPaidPoints(token)
	if err != nil {
		global.Logger.Infof("ошибка на этапе получения информации о выводе стредств %s", err)
		http.Error(writer, err.Error(), statusCode)
		return
	}

	marshal, err := json.Marshal(usersOrder)
	if err != nil {
		global.Logger.Infof("на этапе маршла данных произошла ошибка %s", err)
		http.Error(writer, fmt.Errorf("на этапе маршла данных произошла ошибка %v", err).Error(), 500)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	writer.Write(marshal)
	request.Header.Get("Заказ успешно зарегистрирован и списаны баллы")
}
