package models

type Order struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Value      float64 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
	UserLogin  string
}

type Withdrawn struct {
	ID          string
	Number      string `json:"order"`
	UserLogin   string
	Value       float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
