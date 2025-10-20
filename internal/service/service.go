package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
		rate, err := GetRUBtoUSDRate()
		if err != nil {
			return nil, err
		}
		response.Balance = response.Balance * rate

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
func GetRUBtoUSDRate() (float64, error) {
	key := os.Getenv("EXCHANGERATE_API_KEY")
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/pair/RUB/USD", key)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API вернуло ошибку: %s", resp.Status)
	}

	var data struct {
		ConversionRate float64 `json:"conversion_rate"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.ConversionRate, nil
}
