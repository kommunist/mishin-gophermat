package postwithdrawns

import (
	"encoding/json"
	"io"
	"log/slog"
	"mishin-gophermat/internal/luhn"
	"mishin-gophermat/internal/secure"
	"net/http"
)

type request struct {
	Number string  `json:"order"`
	Value  float64 `json:"sum"`
}

func (h *PostWithdrawsHandler) Process(w http.ResponseWriter, r *http.Request) {
	getLogin := r.Context().Value(secure.UserLoginKey)
	if getLogin == nil {
		slog.Error("Error when get current user from context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	currUser := getLogin.(string)

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
	valid, err := luhn.Valid([]byte(req.Number))
	if err != nil || !valid { // Если не удалось проверить на подлинность, то 422
		slog.Error("Error when valid body", "err", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
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
