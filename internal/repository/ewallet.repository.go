package repository

import (
	"context"
	"errors"
	"log"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EwalletRepository struct {
	db *pgxpool.Pool
}

func NewEWalletRepository(db *pgxpool.Pool) *EwalletRepository {
	return &EwalletRepository{
		db: db,
	}
}

func (er *EwalletRepository) GetBalance(c context.Context, user_id int) (*models.Balance, error) {
	sql := "select user_id, balance from wallets where user_id = $1"

	var balance models.Balance
	if err := er.db.QueryRow(c, sql, user_id).Scan(&balance.User_id, &balance.Balance); err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user_id not found")
		}
		log.Println("Internal Server Error. \nCause: ", err.Error())
		return nil, err
	}
	return &balance, nil
}
