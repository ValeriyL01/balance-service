package models

import "time"

type BalanceResponse struct {
	UserID  int64
	Balance float64
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
