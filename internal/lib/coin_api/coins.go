package coin_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var client = http.DefaultClient

func GetCoins(ctx context.Context, coin string) (int, int64, error) {
	apiURL := fmt.Sprintf("https://min-api.cryptocompare.com/data/pricemulti?fsyms=%s&tsyms=USD", coin)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create request: %w", err)
	}

	data, err := client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to make request: %w", err)
	}
	defer data.Body.Close()

	if data.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("unexpected status code: %d", data.StatusCode)
	}

	body, err := io.ReadAll(data.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read body: %w", err)
	}

	var result map[string]map[string]float64
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to unmarshal body: %w", err)
	}

	dateHeader := data.Header.Get("Date")

	timestamp, err := time.Parse(time.RFC1123, dateHeader)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse date header: %w", err)
	}

	return int(result[coin]["USD"] * 100), timestamp.Unix(), nil
}
