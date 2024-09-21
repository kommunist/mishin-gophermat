package postorders

import (
	"io"
	"log/slog"
	"mishin-gophermat/internal/luhn"
	"net/http"
)

func (h *PostOrdersHandler) Process(w http.ResponseWriter, r *http.Request) {
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

	valid, err := luhn.Valid(body)
	slog.Info("Luhn valid", "valid", valid, "body", string(body))
	if err != nil || !valid { // Если не удалось проверить на подлинность, то 422
		slog.Error("Error when valid body", "err", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	orderData, err := h.DB.OrderByNumberGet(r.Context(), string(body))

	if err != nil {
		slog.Error("Error when find data in db", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if orderData != nil { // если заказ есть
		if orderData["userLogin"] != currUser { // 409 StatusConflict
			slog.Error("Order already upload by another user", "err", err)
			w.WriteHeader(http.StatusConflict)
		} else { // 200 StatusOK
			slog.Info("Order already upload by this user")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		}
	} else { // если заказа нет
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
	}

}
