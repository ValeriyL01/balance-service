package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

func GetUserBalance(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Не валидный id", http.StatusBadRequest)
		return
	}

	response := fmt.Sprintf("User ID: %d", userID)
	_, err = w.Write([]byte(response))
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}
}
