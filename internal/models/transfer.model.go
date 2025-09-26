package models

import "time"

type TransferBody struct {
	IdReceiver    int    `json:"receiver_id" binding:"required"`
	ReceiverPhone int    `json:"receiver_phone" binding:"required"`
	Amount        int    `json:"amount" binding:"required"`
	Notes         string `json:"notes" binding:"required"`
	PinSender     string `json:"pin_sender" binding:"required,min=6"`
}

type UserPin struct {
	Id  int    `db:"id"`
	Pin string `db:"pin"`
}

type TransferResponse struct {
	TransferID     int        `db:"id"`
	TransferStatus string     `json:"transfer_status"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type TransferDatabase struct {
	TransferID     int        `db:"id"`
	SenderWall     int        `db:"sender_wallet_id"`
	ReceiverWall   int        `db:"receiver_wallet_id"`
	Amount         int        `db:"amount"`
	TransferStatus string     `db:"transfer_status"`
	Notes          string     `db:"notes"`
	CreatedAt      *time.Time `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
}
