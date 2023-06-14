package domain

type Order struct {
	ID    int64  `json:"car_id,omitempty"`
	CarID int64  `json:"car_mark,omitempty"`
	Email string `json:"email,omitempty"`
}
