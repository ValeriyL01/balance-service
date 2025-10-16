package service

import (
	"fmt"

	"github.com/ValeriyL01/balance-service/internal/database"
	"github.com/ValeriyL01/balance-service/internal/models"
)

func GetBalance(userID int) (*models.BalanceResponse, error) {
	response, err := database.GetUserBalanceDB(userID)
	if err != nil {
		return nil, fmt.Errorf("service error: %w", err)
	}
	return response, err
}
