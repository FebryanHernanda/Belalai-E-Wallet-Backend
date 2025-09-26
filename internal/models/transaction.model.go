// models/transaction.go
package models

import "time"

type TransactionHistory struct {
	Type           string    `json:"transaction_type" db:"transaction_type"`
	ProfilePicture string    `json:"profile_picture" db:"profile_picture"`
	ContactName    string    `json:"contact_name" db:"contact_name"`
	Amount         string    `json:"amount" db:"display_amount"`
	OriginalAmount int       `json:"original_amount" db:"original_amount"`
	Status         string    `json:"status" db:"status"`
	Notes          string    `json:"notes" db:"notes"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type TransactionHistoryRequest struct {
	UserID int `json:"user_id" db:"user_id"`
}

type TransactionHistoryResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    []TransactionHistory `json:"data,omitempty"`
}
