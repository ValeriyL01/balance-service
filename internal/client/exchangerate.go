package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ValeriyL01/balance-service/internal/utils"
)

func GetRUBtoUSDRate() (float64, error) {
	key := utils.GetEnv("EXCHANGERATE_API_KEY", "")
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
