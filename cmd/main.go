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
		log.Println("Файл .env не найден")
	}

	err = database.ConnectAndInit()
	if err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	defer database.DB.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/balance", handlers.GetUserBalance)
	mux.HandleFunc("/deposit", handlers.DepositBalance)
	mux.HandleFunc("/withdraw", handlers.WithdrawBalance)
	mux.HandleFunc("/transfer", handlers.TransferMoney)
	err = http.ListenAndServe(":4000", mux)

	log.Fatal(err)
}
