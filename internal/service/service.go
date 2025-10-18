package service

import (
	"errors"
	"fmt"

	"github.com/ValeriyL01/balance-service/internal/database"
	"github.com/ValeriyL01/balance-service/internal/models"
)

var ErrInvalidAmount = errors.New("сумма должна быть положительной")
var ErrNoMoney = errors.New("недостаточно средств")

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

func WithdrawBalanceService(balance models.BalanceRequest) error {
	if balance.Amount <= 0 {
		return ErrInvalidAmount
	}
	// что бы списать деньги, сначала нужно получить баланс юзера
	userBalance, err := GetBalance(int(balance.UserID))
	if err != nil {
		return fmt.Errorf("не удалось получить баланс: %w", err)
	}

	if userBalance.Balance <= balance.Amount {
		return ErrNoMoney
	}

	err = database.WithdrawBalanceDB(balance)
	if err != nil {
		return fmt.Errorf("не удалось списать деньги: %w", err)
	}
	return nil

}
