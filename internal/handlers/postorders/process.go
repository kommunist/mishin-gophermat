package postorders

import (
	"io"
	"log/slog"
	"mishin-gophermat/internal/luhn"
	"mishin-gophermat/internal/secure"
	"net/http"
)

func (h *PostOrdersHandler) Process(w http.ResponseWriter, r *http.Request) {
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

	valid, err := luhn.Valid(body)
	slog.Info("Luhn valid", "valid", valid, "body", string(body))
	if err != nil || !valid { // Если не удалось проверить на подлинность, то 422
		slog.Error("Error when valid body", "err", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	orderLogin, err := h.DB.OrderByNumberGet(r.Context(), string(body))

	if err != nil {
		slog.Error("Error when find data in db", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if orderLogin == "" { // если заказа нет
		number := string(body)
		err = h.DB.OrderCreate(r.Context(), number, currUser)
		if err != nil {
			slog.Error("Error when create data in db", "err", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		h.acrChan <- number
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted) // 202 Accepted
		return
	}

	// если заказ есть
	if orderLogin != currUser { // 409 StatusConflict
		slog.Error("Order already upload by another user", "err", err)
		w.WriteHeader(http.StatusConflict)
	} else { // 200 StatusOK
		slog.Info("Order already upload by this user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
