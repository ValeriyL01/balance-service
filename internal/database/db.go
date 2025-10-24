package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ValeriyL01/balance-service/internal/customErrors"
	"github.com/ValeriyL01/balance-service/internal/models"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{db: db}
}

func (d *Database) InitTables() error {
	err := d.createBalancesTable()
	if err != nil {
		return fmt.Errorf("failed to init tables: %w", err)
	}

	err = d.createTransactionTable()
	if err != nil {
		return fmt.Errorf("failed to init tables: %w", err)
	}

	return nil
}

func (d Database) createBalancesTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS balances (
            user_id BIGINT PRIMARY KEY,
            balance DECIMAL(15,2) NOT NULL DEFAULT 0.00,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`

	_, err := d.db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы баланса: %w", err)
	}

	return nil
}

func (d Database) createTransactionTable() error {
	query := `
  CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    type VARCHAR(20) NOT NULL,
    related_user_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)`
	_, err := d.db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы транзакций: %w", err)
	}
	return nil
}

func (d Database) GetUserBalance(userID int) (*models.BalanceResponse, error) {

	userBalance := `SELECT user_id, balance FROM balances WHERE user_id = $1`

	data := d.db.QueryRow(userBalance, userID)

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

func (d Database) DepositBalance(balance models.BalanceRequest) error {

	tx, err := d.db.Begin()
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

	err = d.transactionEntry(balance, "deposit")
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка коммита транзакции: %w", err)
	}

	return nil
}

func (d Database) WithdrawBalance(balance models.BalanceRequest) error {

	tx, err := d.db.Begin()
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

	err = d.transactionEntry(balance, "withdraw")
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка коммита транзакции: %w", err)
	}

	return nil

}

func (d Database) TransferMoney(transfer models.TransferRequest) error {

	tx, err := d.db.Begin()
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

	_, err = d.db.Exec(transactionQuery, transfer.FromUserID, transfer.Amount, "transfer", transfer.ToUserID)
	if err != nil {
		return fmt.Errorf("ошибка записи транзакции: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка коммита транзакции: %w", err)
	}
	return nil
}

// Функция для обновления таблицы транзакций
func (d Database) transactionEntry(balance models.BalanceRequest, tpansactionType string) error {
	transactionQuery := `
INSERT INTO transactions (user_id,amount,type)
VALUES ($1,$2,$3)
`

	_, err := d.db.Exec(transactionQuery, balance.UserID, balance.Amount, tpansactionType)
	if err != nil {
		return fmt.Errorf("ошибка записи транзакции: %w", err)
	}

	return nil
}

func (d Database) GetTransactionUser(userID int, page, limit int, sortBy, sortDir string) (*models.TransactionResponse, error) {

	offset := (page - 1) * limit

	query := fmt.Sprintf(`
        SELECT id, user_id, amount, type, related_user_id, created_at 
        FROM transactions 
        WHERE user_id = $1 
        ORDER BY %s %s 
        LIMIT $2 OFFSET $3
    `, sortBy, sortDir)

	rows, err := d.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction

		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Amount,
			&t.Type,
			&t.RelatedUserID,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения: %w", err)
		}

		transactions = append(transactions, t)
	}
	total, err := d.getTotalTransactions(userID)
	if err != nil {
		return nil, err
	}

	return &models.TransactionResponse{
		Transactions: transactions,
		Total:        total,
		Page:         page,
		PageSize:     limit,
	}, nil
}
func (d Database) getTotalTransactions(userID int) (int, error) {
	var total int
	err := d.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE user_id = $1", userID).Scan(&total)
	return total, err
}
