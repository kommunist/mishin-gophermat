package registration

import (
	"encoding/json"
	"io"
	"log/slog"
	"mishin-gophermat/internal/auth"
	"mishin-gophermat/internal/errors/exist"
	"net/http"
)

type requestItem struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *RegistrationHandler) Process(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPass, err := h.hasher.PassHash(rs.Password)
	if err != nil {
		slog.Error("Error when hash password", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	err = h.DB.UserCreate(r.Context(), rs.Login, hashedPass)

	switch err.(type) { // понравился такой синтаксис проверки, так как можно в любой момент его расширять
	case *exist.ExistError:
		slog.Info("Login already exist")
		w.WriteHeader(http.StatusConflict)
	}

	if err != nil {
		slog.Error("Error when insert user", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	encrypted := auth.Encrypt(map[string]any{"login": rs.Login})
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
