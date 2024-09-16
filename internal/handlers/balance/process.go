package balance

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type response struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func (h *BalanceHandler) Process(w http.ResponseWriter, r *http.Request) {
	var currUser string

	_, claims, _ := h.GetLogin(r.Context())
	if userLogin := claims["login"]; userLogin != nil {
		currUser = claims["login"].(string)
	} else { // 401
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	current, withdrawn, err := h.DB.SelectBalanceByLogin(
		r.Context(), currUser,
	)
	resp := response{Current: current, Withdrawn: withdrawn}

	if err != nil {
		slog.Error("Error when get data from db", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

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
