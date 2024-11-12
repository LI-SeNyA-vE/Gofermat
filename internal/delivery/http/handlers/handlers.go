package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/LI-SeNyA-vE/Gofermat/internal/global"
	"github.com/LI-SeNyA-vE/Gofermat/internal/service"
	"net/http"
)

var log = global.Logger

func UserRegistration(writer http.ResponseWriter, request *http.Request) {
	var (
		buf             bytes.Buffer
		userCredentials service.User
	)

	_, err := buf.ReadFrom(request.Body) //Читает данные из тела запроса
	if err != nil {
		log.Info("")
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &userCredentials) // Разбирает данные из массива byte в структуру
	if err != nil {
		log.Info("")
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = userCredentials.CreateUser() // Функция на создание пользователя
	//Придумать свитч, который будет смотреть если
	//400 — неверный формат запроса;
	//409 — логин уже занят;
	//500 — внутренняя ошибка сервера
	//Пока просто заглушка
	if err != nil {
		log.Info("")
		http.Error(writer, "не найдено", http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}

func UserAuthentication(writer http.ResponseWriter, request *http.Request) {
	var (
		buf             bytes.Buffer
		userCredentials service.User
	)

	_, err := buf.ReadFrom(request.Body) //Читает данные из тела запроса
	if err != nil {
		log.Info("")
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &userCredentials) // Разбирает данные из массива byte в структуру
	if err != nil {
		log.Info("")
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = userCredentials.UserAuthentication() // Функция на создание пользователя
	//Придумать свитч, который будет смотреть если
	//400 — неверный формат запроса;
	//409 — логин уже занят;
	//500 — внутренняя ошибка сервера
	//Пока просто заглушка
	if err != nil {
		log.Info("")
		http.Error(writer, "не найдено", http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}

func AddOrderNumber(writer http.ResponseWriter, request *http.Request) {

}

func ExpenditurePointsOnNewOrder(writer http.ResponseWriter, request *http.Request) {

}

func ListUserOrders(writer http.ResponseWriter, request *http.Request) {

}

func ListUserBalance(writer http.ResponseWriter, request *http.Request) {

}

func InfoAboutUsagePoints(writer http.ResponseWriter, request *http.Request) {

}
