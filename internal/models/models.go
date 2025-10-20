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
	Result             string  `json:"result"`
	Documentation      string  `json:"documentation"`
	TermsOfUse         string  `json:"terms_of_use"`
	TimeLastUpdateUnix int64   `json:"time_last_update_unix"`
	TimeLastUpdateUTC  string  `json:"time_last_update_utc"`
	TimeNextUpdateUnix int64   `json:"time_next_update_unix"`
	TimeNextUpdateUTC  string  `json:"time_next_update_utc"`
	BaseCode           string  `json:"base_code"`
	TargetCode         string  `json:"target_code"`
	ConversionRate     float64 `json:"conversion_rate"`
}
