package balance

import (
	"encoding/json"
	"log/slog"
	"mishin-gophermat/internal/secure"
	"net/http"
)

type response struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func (h *BalanceHandler) Process(w http.ResponseWriter, r *http.Request) {
	getLogin := r.Context().Value(secure.UserLoginKey)
	if getLogin == nil {
		slog.Error("Error when get current user from context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	currUser := getLogin.(string)

	current, withdrawn, err := h.DB.BalanceGet(
		r.Context(), currUser,
	)

	if err != nil {
		slog.Error("Error when get data from db", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp := response{Current: current, Withdrawn: withdrawn}
	respBody, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Error when generate json", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	_, err = w.Write(respBody)
	if err != nil {
		slog.Error("Error when write reponse", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
