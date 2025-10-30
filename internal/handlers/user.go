package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ValeriyL01/balance-service/internal/models"
	"github.com/ValeriyL01/balance-service/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "нужен метод POST", http.StatusMethodNotAllowed)
		return
	}

	var registerData models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerData); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	if registerData.Username == "" || registerData.Email == "" || registerData.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	err := h.userService.Register(registerData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Юзер успешно зарегистрирован",
	})
}
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "нужен метод POST", http.StatusMethodNotAllowed)
		return
	}

	var loginData models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	err := h.userService.Login(loginData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Юзер успешно залогинился",
	})
}
