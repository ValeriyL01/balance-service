package service

import (
	"fmt"

	"github.com/ValeriyL01/balance-service/internal/customErrors"
	"github.com/ValeriyL01/balance-service/internal/database"
	"github.com/ValeriyL01/balance-service/internal/models"
)

func GetBalance(userID int) (*models.BalanceResponse, error) {
	response, err := database.GetUserBalanceDB(userID)
	if err != nil {

		return nil, err
	}
	return response, err
}

func DepositBalanceService(balance models.BalanceRequest) error {
	if balance.Amount <= 0 {
		return customErrors.ErrInvalidAmount
	}

	err := database.DepositBalanceDB(balance)
	if err != nil {
		return fmt.Errorf("не удалось пополнить баланс: %w", err)
	}
	return nil
}

func WithdrawBalanceService(balance models.BalanceRequest) error {
	if balance.Amount <= 0 {
		return customErrors.ErrInvalidAmount
	}

	// что бы списать деньги, сначала нужно получить баланс юзера
	userBalance, err := GetBalance(int(balance.UserID))
	if err != nil {
		return err
	}

	if userBalance.Balance <= balance.Amount {
		return customErrors.ErrNoMoney
	}

	err = database.WithdrawBalanceDB(balance)
	if err != nil {
		return fmt.Errorf("не удалось списать деньги: %w", err)
	}
	return nil

}
