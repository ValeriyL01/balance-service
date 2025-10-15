package main

import (
	"log"
	"net/http"

	"github.com/ValeriyL01/balance-service/internal/database"
	"github.com/ValeriyL01/balance-service/internal/handlers"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден, используем переменные окружения")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal("Ошибка подключения к базе:", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/balance", handlers.GetUserBalance)

	err = http.ListenAndServe(":4000", mux)

	log.Fatal(err)
}
