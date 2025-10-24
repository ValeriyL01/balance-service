package main

import (
	"log"
	"net/http"

	"github.com/ValeriyL01/balance-service/internal/database"
	"github.com/ValeriyL01/balance-service/internal/handlers"
	"github.com/ValeriyL01/balance-service/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	dataBase := database.NewDatabase(db)

	err = dataBase.InitTables()
	if err != nil {
		log.Fatal("Ошибка инициализации таблиц:", err)
	}
	balanceServise := service.NewBalanceService(dataBase)

	handler := handlers.NewHandler(balanceServise)

	mux := http.NewServeMux()
	mux.HandleFunc("/balance", handler.GetUserBalance)
	mux.HandleFunc("/deposit", handler.DepositBalance)
	mux.HandleFunc("/withdraw", handler.WithdrawBalance)
	mux.HandleFunc("/transfer", handler.TransferMoney)
	mux.HandleFunc("/transactions", handler.GetTransactionUser)
	err = http.ListenAndServe(":4000", mux)

	log.Fatal(err)
}
