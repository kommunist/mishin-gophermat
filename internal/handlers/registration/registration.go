package registration

import (
	"encoding/json"
	"io"
	"log/slog"
	"mishin-gophermat/internal/auth"
	"mishin-gophermat/internal/storage"
	"net/http"
)

type requestItem struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Process(w http.ResponseWriter, r *http.Request, db storage.DB) {
	body, err := io.ReadAll(r.Body)
	if err != nil { // если body не учитается
		slog.Error("Error when read body")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	rs := requestItem{}

	err = json.Unmarshal(body, &rs)
	if err != nil || rs.Login == "" || rs.Password == "" { // если запрос не того формата
		slog.Error("Invalid format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = db.CreateUser(r.Context(), rs.Login, rs.Password) // добавить проверку, что существует
	if err != nil {
		slog.Error("Error when insert user")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	encrypted := auth.Encrypt(map[string]interface{}{"login": rs.Login})
	newCookie := newAuthCookie(encrypted)
	http.SetCookie(w, &newCookie)
	w.Header().Set("Authorization", "BEARER "+encrypted)
}

func newAuthCookie(value string) http.Cookie {
	return http.Cookie{
		Name:  "jwt",
		Value: value,
		Path:  "/",
	}
}
