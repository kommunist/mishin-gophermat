package login

import (
	"encoding/json"
	"io"
	"log/slog"
	"mishin-gophermat/internal/secure"
	"net/http"
)

type requestItem struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *LoginHandler) Process(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil { // если body не учитается
		slog.Error("Error when read body", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	rs := requestItem{}

	err = json.Unmarshal(body, &rs)
	if err != nil || rs.Login == "" || rs.Password == "" { // если запрос не того формата
		slog.Error("Invalid format", "err", err)
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	hashed, err := h.DB.UserGet(r.Context(), rs.Login)
	if err != nil {
		slog.Error("Error when select password of user", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	if !h.checker.PassCheck(rs.Password, hashed) {
		slog.Info("User not found")
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}

	encrypted := secure.EncryptLogin(rs.Login)
	newCookie := newAuthCookie(encrypted)
	http.SetCookie(w, &newCookie)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "BEARER "+encrypted)
}

func newAuthCookie(value string) http.Cookie {
	return http.Cookie{
		Name:  "jwt",
		Value: value,
		Path:  "/",
	}
}
