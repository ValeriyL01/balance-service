package database

import (
	"database/sql"
	"fmt"

	"github.com/ValeriyL01/balance-service/internal/utils"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	dbUser := utils.GetEnv("DB_USER", "")
	dbPassword := utils.GetEnv("DB_PASSWORD", "")
	dbName := utils.GetEnv("DB_NAME", "")
	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbSSLMode := utils.GetEnv("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		dbUser, dbPassword, dbName, dbHost, dbPort, dbSSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
