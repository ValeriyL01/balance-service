package service

import (
	"fmt"
	"math"

	"github.com/ValeriyL01/balance-service/internal/api"
	"github.com/ValeriyL01/balance-service/internal/customErrors"
	"github.com/ValeriyL01/balance-service/internal/database"

	"github.com/ValeriyL01/balance-service/internal/models"
)

func GetBalance(userID int, currency string) (*models.BalanceResponse, error) {
	response, err := database.GetUserBalanceDB(userID)

	if err != nil {

		return nil, err
	}

	if currency == "USD" {
		rate, err := api.GetRUBtoUSDRate()
		if err != nil {
			return nil, err
		}
		response.Balance = response.Balance * rate

		response.Balance = math.Round(response.Balance*100) / 100
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
	userBalance, err := GetBalance(int(balance.UserID), "")
	if err != nil {
		return err
	}

	if userBalance.Balance < balance.Amount {
		return customErrors.ErrNoMoney
	}

	err = database.WithdrawBalanceDB(balance)
	if err != nil {
		return fmt.Errorf("не удалось списать деньги: %w", err)
	}
	return nil

}

func TransferMoneyService(transfer models.TransferRequest) error {
	if transfer.Amount <= 0 {
		return customErrors.ErrInvalidAmount
	}

	userBalance, err := GetBalance(int(transfer.FromUserID), "")
	if err != nil {

		return err
	}
	_, err = GetBalance(int(transfer.ToUserID), "")
	if err != nil {

		return err
	}
	if userBalance.Balance < transfer.Amount {
		return customErrors.ErrNoMoney
	}

	err = database.TransferMoneyDB(transfer)
	if err != nil {
		return fmt.Errorf("не удалось перевести деньги: %w", err)
	}
	return nil

}

func GetTransactionUserService(userID, page, limit int, sortBy, sortDir string) (*models.TransactionResponse, error) {
	_, err := GetBalance(userID, "")
	if err != nil {

		return nil, err
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	if sortBy != "amount" {
		sortBy = "created_at"
	}
	if sortDir != "asc" {
		sortDir = "desc"
	}
	response := &models.TransactionResponse{}

	response, err = database.GetTransactionUserDB(userID, page, limit, sortBy, sortDir)
	if err != nil {
		return nil, err
	}
	return response, nil
}
