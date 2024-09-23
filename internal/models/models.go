package models

type Order struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Value      float64 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
	UserLogin  string
}
