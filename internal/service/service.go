package service

import (
	"fmt"
	"math"

	"github.com/ValeriyL01/balance-service/internal/api"
	"github.com/ValeriyL01/balance-service/internal/customErrors"
	"github.com/ValeriyL01/balance-service/internal/database"

	"github.com/ValeriyL01/balance-service/internal/models"
)

type BalanceService struct {
	dataBaseStruct *database.Database
}

func NewBalanceService(dataBaseStruct *database.Database) *BalanceService {
	return &BalanceService{dataBaseStruct: dataBaseStruct}
}

func (b BalanceService) GetBalance(userID int, currency string) (*models.BalanceResponse, error) {
	response, err := b.dataBaseStruct.GetUserBalance(userID)

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

func (b BalanceService) DepositBalance(balance models.BalanceRequest) error {
	if balance.Amount <= 0 {
		return customErrors.ErrInvalidAmount
	}

	err := b.dataBaseStruct.DepositBalance(balance)
	if err != nil {
		return fmt.Errorf("не удалось пополнить баланс: %w", err)
	}
	return nil
}

func (b BalanceService) WithdrawBalance(balance models.BalanceRequest) error {
	if balance.Amount <= 0 {
		return customErrors.ErrInvalidAmount
	}

	// что бы списать деньги, сначала нужно получить баланс юзера
	userBalance, err := b.GetBalance(int(balance.UserID), "")
	if err != nil {
		return err
	}

	if userBalance.Balance < balance.Amount {
		return customErrors.ErrNoMoney
	}

	err = b.dataBaseStruct.WithdrawBalance(balance)
	if err != nil {
		return fmt.Errorf("не удалось списать деньги: %w", err)
	}
	return nil

}

func (b BalanceService) TransferMoney(transfer models.TransferRequest) error {
	if transfer.Amount <= 0 {
		return customErrors.ErrInvalidAmount
	}

	userBalance, err := b.GetBalance(int(transfer.FromUserID), "")
	if err != nil {

		return err
	}
	_, err = b.GetBalance(int(transfer.ToUserID), "")
	if err != nil {

		return err
	}
	if userBalance.Balance < transfer.Amount {
		return customErrors.ErrNoMoney
	}

	err = b.dataBaseStruct.TransferMoney(transfer)
	if err != nil {
		return fmt.Errorf("не удалось перевести деньги: %w", err)
	}
	return nil

}

func (b BalanceService) GetTransactionUser(userID, page, limit int, sortBy, sortDir string) (*models.TransactionResponse, error) {
	_, err := b.GetBalance(userID, "")
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

	response, err = b.dataBaseStruct.GetTransactionUser(userID, page, limit, sortBy, sortDir)
	if err != nil {
		return nil, err
	}
	return response, nil
}
