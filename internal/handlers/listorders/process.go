package listorders

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// type responseItem struct {
// 	Number     string  `json:"number"`
// 	Status     string  `json:"status"`
// 	Accrual    float64 `json:"accrual"`
// 	UploadedAt string  `json:"uploaded_at"`
// }

func (h *ListOrdersHandler) Process(w http.ResponseWriter, r *http.Request) {
	var currUser string
	// resp := make([]responseItem, 0)

	_, claims, _ := h.GetLogin(r.Context())
	if userLogin := claims["login"]; userLogin != nil {
		currUser = claims["login"].(string)
	} else { // 401
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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

	// for _, v := range data {
	// 	resp = append(
	// 		resp,
	// 		responseItem{
	// 			Number:     v["number"].(string),
	// 			Status:     v["status"].(string),
	// 			Accrual:    v["accrual"].(float64),
	// 			UploadedAt: v["uploadedAt"].(string),
	// 		},
	// 	)
	// }

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
