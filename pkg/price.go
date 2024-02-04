package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PriceResponse struct {
	ID     string  `json:"id"`
	Symbol string  `json:"symbol"`
	Name   string  `json:"name"`
	Price  float64 `json:"current_price"`
}

func GetCoinPrice(coinID string) (float64, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/markets?vs_currency=USD&ids=%s&order=market_cap_desc&per_page=1&page=1&sparkline=false", coinID)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var response []PriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}

	if len(response) == 0 {
		return 0, fmt.Errorf("Coin not found")
	}

	for _, coin := range response {
		if coin.ID == coinID {
			return coin.Price, nil
		}
	}

	return 0, fmt.Errorf("Coin not found")
}
