package models

import "time"

type BalanceResponse struct {
	UserID  int64   `json:"user_id"`
	Balance float64 `json:"balance"`
}
type Balance struct {
	UserID    int64
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Transaction struct {
	ID            int64
	UserID        int64
	Amount        float64
	Type          string
	RelatedUserID *int64
	CreatedAt     time.Time
}

type TransactionResponse struct {
	Transactions []Transaction `json:"transactions"`
	Total        int           `json:"total"`
	Page         int           `json:"page"`
	PageSize     int           `json:"page_size"`
}

type BalanceRequest struct {
	UserID int64   `json:"user_id"`
	Amount float64 `json:"amount"`
}

type TransferRequest struct {
	FromUserID int64   `json:"from_user_id"`
	ToUserID   int64   `json:"to_user_id"`
	Amount     float64 `json:"amount"`
}

// структура для ответа от api по обмену валют для вывода счета в долларах
type ExchangeRateResponse struct {
	ConversionRate float64 `json:"conversion_rate"`
}
