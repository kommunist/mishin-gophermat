package postwithdrawns

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type request struct {
	Number string  `json:"order"`
	Value  float64 `json:"sum"`
}

// не стал реализовывать кейс с неправильным номером заказа

func (h *PostWithdrawsHandler) Process(w http.ResponseWriter, r *http.Request) {
	var currUser string

	_, claims, _ := h.GetLogin(r.Context())
	if userLogin := claims["login"]; userLogin != nil {
		currUser = claims["login"].(string)
	} else { // 401
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil { // если body не читается
		slog.Error("Error when read body", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	req := request{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		slog.Error("Error when parse json", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// проверим баланс
	curr, _, err := h.DB.BalanceGet(r.Context(), currUser)
	if err != nil {
		slog.Error("Error when select balance of current user", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if req.Value > curr {
		slog.Info("Balance is low")
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}

	err = h.DB.WithdrawnCreate(r.Context(), currUser, req.Number, req.Value)
	if err != nil {
		slog.Error("Error when create withdrawn in db", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 Accepted
}
