package registration

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type requestItem struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Process(w http.ResponseWriter, r *http.Request, db string) {
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
}
