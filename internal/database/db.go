package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/ValeriyL01/balance-service/internal/models"
	_ "github.com/lib/pq"
)

// Глобальная переменная для подключения к БД
var DB *sql.DB

func ConnectAndInit() error {
	db, err := Connect()
	if err != nil {
		return err
	}

	// Инициализируем глобальную переменную
	DB = db

	err = createBalancesTable()

	if err != nil {
		db.Close()
		return fmt.Errorf("failed to init tables: %w", err)
	}
	err = createTransactionTable()
	if err != nil {
		db.Close()
		return fmt.Errorf("failed to init tables: %w", err)
	}
	return nil
}

func Connect() (*sql.DB, error) {
	dbUser := getEnv("DB_USER", "")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func createBalancesTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS balances (
            user_id BIGINT PRIMARY KEY,
            balance DECIMAL(15,2) NOT NULL DEFAULT 0.00,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы баланса: %w", err)
	}

	return nil
}

func createTransactionTable() error {
	query := `
  CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    type VARCHAR(20) NOT NULL,
    related_user_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)`
	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы транзакций: %w", err)
	}
	return nil
}

func GetUserBalanceDB(userID int) (*models.BalanceResponse, error) {

	if DB == nil {
		return nil, fmt.Errorf("база данных не инициализирована ")
	}

	userBalance := `SELECT user_id, balance FROM balances WHERE user_id = $1`

	data := DB.QueryRow(userBalance, userID)

	balance := &models.BalanceResponse{}

	err := data.Scan(&balance.UserID, &balance.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user c id: %d не найден", userID)
		}
		return nil, fmt.Errorf(" баланс не получен для юзера c id: %d: %w", userID, err)
	}

	return balance, nil
}

func DepositBalanceDB(balance models.BalanceRequest) error {
	if DB == nil {
		return fmt.Errorf("база данных не инициализирована ")
	}
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("ошибка транзакции: %w", err)
	}
	defer tx.Rollback()
	balanceQuery := `
		INSERT INTO balances (user_id, balance) 
		VALUES ($1, $2)
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			balance = balances.balance + $2,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err = tx.Exec(balanceQuery, balance.UserID, balance.Amount)
	if err != nil {
		return fmt.Errorf("ошибка обновления баланса: %w", err)
	}

	transactionQuery := `
INSERT INTO transactions (user_id,amount,type)
VALUES ($1,$2,$3)
`

	_, err = tx.Exec(transactionQuery, balance.UserID, balance.Amount, "deposit")
	if err != nil {
		return fmt.Errorf("ошибка записи транзакции: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка коммита транзакции: %w", err)
	}

	return nil
}
