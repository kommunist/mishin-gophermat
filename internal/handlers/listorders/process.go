package listorders

import (
	"encoding/json"
	"log/slog"
	"mishin-gophermat/internal/secure"
	"net/http"
)

func (h *ListOrdersHandler) Process(w http.ResponseWriter, r *http.Request) {
	getLogin := r.Context().Value(secure.UserLoginKey)
	if getLogin == nil {
		slog.Error("Error when get current user from context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	currUser := getLogin.(string)

	data, err := h.DB.OrdersGet(r.Context(), currUser)
	if err != nil { // 204
		slog.Error("Error when get data from db", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if len(data) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return

	}

	respBody, err := json.Marshal(data)
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
