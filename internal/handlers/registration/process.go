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

	err = h.DB.CreateUser(r.Context(), rs.Login, rs.Password) // добавить проверку, что существует

	switch err.(type) { // понравился такой синтаксис проверки, так как можно в любой момент его расширять
	case *exist.ExistError:
		slog.Info("Login already exist")
		w.WriteHeader(http.StatusConflict)
	}

	if err != nil {
		slog.Error("Error when insert user", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	encrypted := auth.Encrypt(map[string]interface{}{"login": rs.Login})
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
