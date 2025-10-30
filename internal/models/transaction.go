package models

import "time"

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

// структура для ответа от api по обмену валют для вывода счета в долларах
type ExchangeRateResponse struct {
	ConversionRate float64 `json:"conversion_rate"`
}
