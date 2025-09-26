package models

import "time"

type TransferBody struct {
	IdReceiver    int    `json:"receiver_id"`
	ReceiverPhone int    `json:"receiver_phone"`
	Amount        int    `json:"amount"`
	Notes         string `json:"notes"`
	PinSender     int    `json:"pin_sender"`
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
