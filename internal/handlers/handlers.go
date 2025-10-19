package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/ValeriyL01/balance-service/internal/customErrors"
	"github.com/ValeriyL01/balance-service/internal/models"
	"github.com/ValeriyL01/balance-service/internal/service"
)

func GetUserBalance(w http.ResponseWriter, r *http.Request) {

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Укажите user_id", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "user_id должен быть числом", http.StatusBadRequest)
		return
	}

	response, err := service.GetBalance(userID)
	if err != nil {

		if errors.Is(err, customErrors.ErrUserNotFound) {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}

		http.Error(w, "Не удалось получить баланс", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DepositBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "нужен метод POST", http.StatusMethodNotAllowed)
		return
	}
	var request models.BalanceRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Ошибка парсинга JSON: %v", err)
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	err := service.DepositBalanceService(request)
	if err != nil {
		log.Printf("Ошибка при пополнении баланса: %v", err)

		if errors.Is(err, customErrors.ErrInvalidAmount) {
			http.Error(w, "Сумма должна быть положительной", http.StatusBadRequest)
		} else {
			http.Error(w, "Не удалось пополнить баланс", http.StatusInternalServerError)
			return
		}

	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Баланс успешно пополнен",
	})

}

func WithdrawBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "нужен метод POST", http.StatusMethodNotAllowed)
	}

	var request models.BalanceRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Ошибка парсинга JSON: %v", err)
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	err := service.WithdrawBalanceService(request)
	if err != nil {

		if errors.Is(err, customErrors.ErrInvalidAmount) {
			http.Error(w, "Сумма должна быть положительной", http.StatusBadRequest)
			return
		} else if errors.Is(err, customErrors.ErrNoMoney) {
			http.Error(w, "Недостаточно средств", http.StatusBadRequest)
			return
		} else if errors.Is(err, customErrors.ErrUserNotFound) {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		} else {
			http.Error(w, "списание не удалось", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Успешное списание",
	})
}
