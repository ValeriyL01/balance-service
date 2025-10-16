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
