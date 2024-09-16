package listwithdrawns

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type responseItem struct {
	Number      string `json:"order"`
	Value       int    `json:"sum"`
	ProcessedAt string `json:"processed_at"`
}

func (h *ListWithdrawns) Process(w http.ResponseWriter, r *http.Request) {
	var currUser string
	resp := make([]responseItem, 0)

	_, claims, _ := h.GetLogin(r.Context())
	if userLogin := claims["login"]; userLogin != nil {
		currUser = claims["login"].(string)
	} else { // 401
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := h.DB.SelectWithdrawnsByLogin(r.Context(), currUser)
	if err != nil {
		slog.Error("Error when get data from db", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if len(data) == 0 {
		w.WriteHeader(http.StatusNoContent) // 204
		return

	}
	for _, v := range data {
		resp = append(
			resp,
			responseItem{
				Number:      v["number"].(string),
				Value:       v["value"].(int),
				ProcessedAt: v["processedAt"].(string),
			},
		)
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
