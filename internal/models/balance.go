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
type BalanceRequest struct {
	UserID int64   `json:"user_id"`
	Amount float64 `json:"amount"`
}
type TransferRequest struct {
	FromUserID int64   `json:"from_user_id"`
	ToUserID   int64   `json:"to_user_id"`
	Amount     float64 `json:"amount"`
}
