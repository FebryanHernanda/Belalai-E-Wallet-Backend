package models

import "time"

type TopUpStatus string

const (
	TopUpSuccess TopUpStatus = "success"
	TopUpFailed  TopUpStatus = "failed"
	TopUpPending TopUpStatus = "pending"
)

type PaymentMethod struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type TopUp struct {
	ID        int         `db:"id" json:"id"`
	Amount    int         `db:"amount" json:"amount"`
	Tax       int         `db:"tax" json:"tax"`
	PaymentID int         `db:"payment_id" json:"payment_id"`
	Status    TopUpStatus `db:"topup_status" json:"topup_status"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time  `db:"updated_at" json:"updated_at,omitempty"`
}

type TopUpRequest struct {
	Amount    int `json:"amount" binding:"required"`
	Tax       int `json:"tax"`
	PaymentID int `json:"payment_id" binding:"required"`
}

type TopUpResponse struct {
	ID        int         `json:"id"`
	Amount    int         `json:"amount"`
	Tax       int         `json:"tax"`
	PaymentID int         `json:"payment_id"`
	Status    TopUpStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
}
