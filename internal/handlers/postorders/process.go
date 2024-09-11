package postorders

import (
	"fmt"
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
		slog.Error("Error when read body")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	valid, err := luhn.Valid(body)
	slog.Info("Luhn valid", "valid", valid, "body", string(body))
	if err != nil || !valid { // Если не удалось проверить на подлинность, то 422
		slog.Error("Error when valid body")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	orderData, err := h.DB.SelectOrderByNumber(r.Context(), string(body))

	if err != nil {
		slog.Error("Error when find data in db")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if orderData != nil { // если заказ есть
		if orderData["userLogin"] != currUser { // 409 StatusConflict
			slog.Error("Order already upload by another user")
			w.WriteHeader(http.StatusConflict)
		} else { // 200 StatusOK
			slog.Info("Order already upload by this user")
			w.WriteHeader(http.StatusOK)
		}
	} else { // если заказа нет
		fmt.Println(currUser)
		err = h.DB.CreateOrder(r.Context(), string(body), currUser)
		if err != nil {
			slog.Error("Error when create data in db")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted) // 202 Accepted
	}

}
