package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
		http.Error(w, "Не удалось получить баланс", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
