package service

import (
	"errors"
	"fmt"

	"github.com/ValeriyL01/balance-service/internal/database"
	"github.com/ValeriyL01/balance-service/internal/models"
)

var ErrInvalidAmount = errors.New("сумма должна быть положительной")

func GetBalance(userID int) (*models.BalanceResponse, error) {
	response, err := database.GetUserBalanceDB(userID)
	if err != nil {
		return nil, fmt.Errorf("service error: %w", err)
	}
	return response, err
}

func DepositBalanceService(balance models.BalanceRequest) error {
	if balance.Amount <= 0 {
		return ErrInvalidAmount
	}

	err := database.DepositBalanceDB(balance)
	if err != nil {
		return fmt.Errorf("не удалось пополнить баланс: %w", err)
	}
	return nil
}
