package repository

import (
	"context"
	"errors"
	"log"
	"math"
	"time"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/Belalai-E-Wallet-Backend/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type TransferRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewTransferRepository(db *pgxpool.Pool, rdb *redis.Client) *TransferRepository {
	return &TransferRepository{db: db, rdb: rdb}
}

// filter user by name or phone number
func (ur *TransferRepository) FilterUser(c context.Context, query string, offset, limit, page int) (models.ListprofileResponse, error) {

	// Inisialisasi response
	searchQuery := "%" + query + "%"
	values := []any{searchQuery, limit, offset}

	// Get cached filter user, before accesing database
	// Get and renew cache only for page 1 only (offset 0)
	rdbKey := "Belalai-E-wallet:filter-user"
	if offset == 0 && query == "" {
		cachedUser, err := utils.RedisGetData[models.ListprofileResponse](c, *ur.rdb, rdbKey)
		if err != nil {
			log.Println("Redis error :", err)
		} else if cachedUser != nil && len(cachedUser.Users) > 0 {
			return *cachedUser, nil
		}
	}

	// if there is no key/ no cached data
	// get filter user from database
	// get total filtered user
	countSql := `SELECT COUNT(user_id) FROM profile
    WHERE fullname ILIKE $1 OR phone ILIKE $1`
	var totalUser int
	countErr := ur.db.QueryRow(c, countSql, searchQuery).Scan(&totalUser)
	if countErr != nil {
		// Tangani error, bisa jadi error umum atau NoRows, tapi COUNT() biasanya selalu mengembalikan 0 atau lebih.
		log.Println("Error getting total user count:", countErr.Error())
		return models.ListprofileResponse{}, countErr
	}

	// get filtered list user
	sql := `SELECT user_id, profile_picture, fullname, phone FROM profile
    WHERE fullname ILIKE $1 OR phone ILIKE $1
		LIMIT $2 OFFSET $3`
	rows, err := ur.db.Query(c, sql, values...)
	if err != nil {
		log.Println("internal server error : ", err.Error())
		return models.ListprofileResponse{}, err
	}
	defer rows.Close()

	// processing data / read rows
	var users []models.ProfileResponse
	for rows.Next() {
		var user models.ProfileResponse
		if err := rows.Scan(&user.UserID, &user.ProfilePicture, &user.Fullname, &user.Phone); err != nil {
			log.Println("Scan Error, ", err.Error())
			return models.ListprofileResponse{}, err
		}
		users = append(users, user)
	}

	// final response
	finalResponse := models.ListprofileResponse{
		Users:     users,
		Page:      page,
		Limit:     limit,
		TotalUser: totalUser,
		TotalPage: int(math.Ceil(float64(totalUser) / float64(limit))),
	}

	// make cache filter user after query data from database
	// Get and renew cache only for page 1 only (offset 0)
	if offset == 0 && query == "" {
		if err := utils.RedisRenewData(c, *ur.rdb, rdbKey, finalResponse, 10*time.Minute); err != nil {
			log.Println("Failed to renew Redis cache:", err.Error())
		}
	}

	// return users, and error nil if success
	return finalResponse, nil
}

// get user pin for validate user on handler
func (ur *TransferRepository) GetHashedPin(rqCntxt context.Context, senderId int) (models.UserPin, error) {
	var userPin models.UserPin
	sql := `SELECT id, pin FROM users WHERE id=$1`
	if err := ur.db.QueryRow(rqCntxt, sql, senderId).Scan(&userPin.Id, &userPin.Pin); err != nil {
		log.Println("failed get pin user \ncause :", err)
		return models.UserPin{}, errors.New("failed get pin user")
	}
	return userPin, nil
}

// transfer transaction
var ErrNotEnoughBalance = errors.New("not enough balance for this transfer")
var ErrCantSendingToYourself = errors.New("can't sending money to yourself")

func (ur *TransferRepository) TransferMoney(rqCntxt context.Context, senderId int, body models.TransferBody) error {

	// using tx transaction postgresql
	tx, err := ur.db.Begin(rqCntxt)
	if err != nil {
		log.Println("Failed to begin DB transaction\nCause: ", err)
		return err
	}
	defer tx.Rollback(rqCntxt)

	// get balance sender and validate pin sender
	// if balance sender is not enough to do transfer, abort transaction
	var senderWalletID int
	var senderBalance float64
	qBalSender := `SELECT w.id, w.balance FROM wallets w 
         JOIN users u ON w.user_id = u.id 
         WHERE u.id = $1 FOR UPDATE`
	if err := tx.QueryRow(rqCntxt, qBalSender, senderId).Scan(&senderWalletID, &senderBalance); err != nil {
		if err == pgx.ErrNoRows {
			log.Println("error no rows or user invalid", err)
			return errors.New("user invalid or wrong pin inputs")
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		return err
	}
	// validate not sending money to self
	if senderWalletID == body.IdReceiver {
		return ErrCantSendingToYourself
	}

	// validate if sender balance is have enough money to do transfer
	if senderBalance < float64(body.Amount) {
		return ErrNotEnoughBalance
	}

	// execute tranfser
	// update saldo sender
	now := time.Now()
	sqlSenderWallet := `UPDATE wallets SET balance = balance - $1, updated_at = $2 WHERE id = $3`
	values := []any{body.Amount, now, senderId}
	cmd, err := tx.Exec(rqCntxt, sqlSenderWallet, values...)
	if err != nil {
		log.Println("Failed execute query sqlSenderWallet\nCause:", err)
		return err
	}
	if cmd.RowsAffected() == 0 {
		log.Println("no row effected when UPDATE wallets maybe failed?")
		return errors.New("no row effected when UPDATE wallets maybe failed?")
	}
	// update saldo receiver
	sqlReceiverWallet := `UPDATE wallets SET balance = balance + $1, updated_at = $2 WHERE id = $3`
	values = []any{body.Amount, now, body.IdReceiver}
	cmd, err = tx.Exec(rqCntxt, sqlReceiverWallet, values...)
	if err != nil {
		log.Println("Failed execute query sqlReceiverWallet\nCause:", err)
		return err
	}
	if cmd.RowsAffected() == 0 {
		log.Println("no row effected when UPDATE wallets maybe failed?")
		return errors.New("no row effected when UPDATE wallets maybe failed?")
	}

	// insert transfer data
	var transferID int
	sqlTansferTable := `INSERT INTO transfer (sender_wallet_id, receiver_wallet_id, amount, transfer_status, notes, created_at, updated_at)
    VALUES ($1, $2, $3, 'success', $4, $5, $5) RETURNING id`
	values = []any{senderId, body.IdReceiver, body.Amount, body.Notes, now}
	if err := tx.QueryRow(rqCntxt, sqlTansferTable, values...).Scan(&transferID); err != nil {
		log.Println("Failed execute query sqlTansferTable \nCause :", err)
		return err
	}

	// insert wallet_transfer
	sqlTransferWalletTable := `INSERT INTO wallets_transfer (wallets_id, transfer_id) VALUES ($1, $3), ($2, $3)`
	values = []any{senderId, body.IdReceiver, transferID}
	cmd, err = tx.Exec(rqCntxt, sqlTransferWalletTable, values...)
	if err != nil {
		log.Println("Failed execute query sqlTransferWalletTable\nCause:", err)
		return err
	}
	if cmd.RowsAffected() == 0 {
		log.Println("no row effected when INSERT INTO wallets_transfer maybe failed?")
		return errors.New("no row effected when INSERT INTO wallets_transfer maybe failed?")
	}

	// commit transaction if all query success execute
	if err := tx.Commit(rqCntxt); err != nil {
		log.Println("Failed to commit DB transaction\nCause: ", err)
		return err
	}
	log.Println("success to commit DB transaction")

	// error nil if success
	return nil
}
