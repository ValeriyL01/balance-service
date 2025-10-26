package database

import (
	"database/sql"
	"fmt"

	"github.com/ValeriyL01/balance-service/internal/config"
	_ "github.com/lib/pq"
)

func Connect(config config.DB) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		config.User, config.Password, config.Name, config.Host, config.Port, config.SSLMode)

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
