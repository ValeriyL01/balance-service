package main

import (
	"log"

	"github.com/ValeriyL01/balance-service/internal/config"
	"github.com/ValeriyL01/balance-service/internal/database"
	"github.com/ValeriyL01/balance-service/internal/handlers"
	"github.com/ValeriyL01/balance-service/internal/server"
	"github.com/ValeriyL01/balance-service/internal/service"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Parse()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg.DB)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	dataBase := database.NewDatabase(db)
	userDB := database.NewUserDB(db)
	err = dataBase.InitTables()
	if err != nil {
		log.Fatal("Ошибка инициализации таблиц:", err)
	}
	balanceService := service.NewBalanceService(dataBase)
	userService := service.NewUserService(userDB)

	handler := handlers.NewHandler(balanceService)
	userHandler := handlers.NewUserHandler(userService)
	srv := server.NewServer(cfg.Port, handler, userHandler)

	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
