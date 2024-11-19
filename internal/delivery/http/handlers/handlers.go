package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/LI-SeNyA-vE/Gofermat/internal/global"
	"github.com/LI-SeNyA-vE/Gofermat/internal/service"
	"net/http"
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

func AddOrderNumber(writer http.ResponseWriter, request *http.Request) {
	var (
		buf         bytes.Buffer
		numberOrder int
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

	err = json.Unmarshal(buf.Bytes(), &numberOrder) // Разбирает данные из массива byte в структуру
	if err != nil {
		global.Logger.Info("в заказе присутствуют другие символы кроме цифр")
		http.Error(writer, fmt.Sprintf("в заказе присутствуют другие символы кроме цифр %v", err), http.StatusUnprocessableEntity)
		return
	}

	statusCode, err := service.UserUploadingNumberOrder(request.Header.Get("Authorization"), numberOrder)
	if err != nil {
		global.Logger.Infof("ошибка на этапе загрузки пользователем заказа %s", err)
		http.Error(writer, err.Error(), statusCode)
		return
	}
	writer.WriteHeader(statusCode)
}

func ExpenditurePointsOnNewOrder(writer http.ResponseWriter, request *http.Request) {

}

func ListUserOrders(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("Content-Length") != "0" {
		http.Error(writer, "Content-Length != 0", http.StatusBadRequest)
		return
	}
	ordersJSON, statusCode, err := service.UserListUserOrders(request.Header.Get("Authorization"))
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
	balanceJSON, statusCode, err := service.UserListUserBalance(request.Header.Get("Authorization"))
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

}
