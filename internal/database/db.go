package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ValeriyL01/balance-service/internal/customErrors"
	"github.com/ValeriyL01/balance-service/internal/models"
	"github.com/ValeriyL01/balance-service/internal/utils"

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
	dbUser := utils.GetEnv("DB_USER", "")
	dbPassword := utils.GetEnv("DB_PASSWORD", "")
	dbName := utils.GetEnv("DB_NAME", "")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=localhost port=5432 sslmode=disable",
		dbUser, dbPassword, dbName)

	// создает объект *sql.DB для работы с базой
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// подключение к базе данных
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
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

	userBalance := `SELECT user_id, balance FROM balances WHERE user_id = $1`

	data := DB.QueryRow(userBalance, userID)

	balance := &models.BalanceResponse{}

	err := data.Scan(&balance.UserID, &balance.Balance)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {

			return nil, customErrors.ErrUserNotFound
		}
		return nil, fmt.Errorf(" баланс не получен для юзера c id: %d: %w", userID, err)
	}

	return balance, nil
}

func DepositBalanceDB(balance models.BalanceRequest) error {

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

	err = transactionEntry(balance, "deposit")
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка коммита транзакции: %w", err)
	}

	return nil
}

func WithdrawBalanceDB(balance models.BalanceRequest) error {

	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("ошибка транзакции: %w", err)
	}
	defer tx.Rollback()

	balanceQuery := `
UPDATE balances 
SET balance = balance - $1,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $2 
`

	_, err = tx.Exec(balanceQuery, balance.Amount, balance.UserID)
	if err != nil {
		return fmt.Errorf("ошибка обновления баланса: %w", err)
	}

	err = transactionEntry(balance, "withdraw")
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка коммита транзакции: %w", err)
	}

	return nil

}

func TransferMoneyDB(transfer models.TransferRequest) error {

	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("ошибка транзакции: %w", err)
	}
	defer tx.Rollback()
	transferFromQuery := `
UPDATE balances 
SET balance = balance - $1,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $2 
`

	_, err = tx.Exec(transferFromQuery, transfer.Amount, transfer.FromUserID)
	if err != nil {
		return fmt.Errorf("ошибка обновления баланса: %w", err)
	}
	transferToFromQuery := `
	UPDATE balances 
SET balance = balance + $1,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $2 
	`

	_, err = tx.Exec(transferToFromQuery, transfer.Amount, transfer.ToUserID)
	if err != nil {
		return fmt.Errorf("ошибка обновления баланса: %w", err)
	}

	transactionQuery := `
INSERT INTO transactions (user_id,amount,type,related_user_id)
VALUES ($1,$2,$3,$4)
`

	_, err = DB.Exec(transactionQuery, transfer.FromUserID, transfer.Amount, "transfer", transfer.ToUserID)
	if err != nil {
		return fmt.Errorf("ошибка записи транзакции: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка коммита транзакции: %w", err)
	}
	return nil
}

// Функция для обновления таблицы транзакций
func transactionEntry(balance models.BalanceRequest, tpansactionType string) error {
	transactionQuery := `
INSERT INTO transactions (user_id,amount,type)
VALUES ($1,$2,$3)
`

	_, err := DB.Exec(transactionQuery, balance.UserID, balance.Amount, tpansactionType)
	if err != nil {
		return fmt.Errorf("ошибка записи транзакции: %w", err)
	}

	return nil
}
