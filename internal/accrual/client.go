package accrual

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type response struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func (acr *Accrual) getOrderData(number string) (string, float64, int, error) { // status, accrual, wait, error
	// req = send_request
	resp, err := http.Get(acr.URI + "/api/orders/" + number)
	if err != nil {
		slog.Error("Error when send request to accruals", "err", err)
		return "", 0, 0, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		wait, err := strconv.Atoi(resp.Header.Get("Retry-After"))
		if err != nil {
			slog.Error("Error when convert data fron 429 status", err, err)
			return "", 0, 0, err
		}
		return "", 0, wait, nil
	}

	if resp.StatusCode != http.StatusOK {
		slog.Info("Accrual return somthing another", "statusCode", resp.StatusCode)
		return "", 0, 0, nil
	}

	bytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		slog.Error("Error when get bytes from response", "err", err)
		return "", 0, 0, err
	}
	respStr := response{}
	err = json.Unmarshal(bytes, &respStr)
	if err != nil {
		slog.Error("Error when parsing json", "err", err)
		return "", 0, 0, err
	}

	return respStr.Status, respStr.Accrual, 0, nil
}
