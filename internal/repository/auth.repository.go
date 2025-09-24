package repository

import (
	"context"
	"errors"
	"log"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (ar *AuthRepository) GetEmail(c context.Context, email string) (*models.User, error) {
	sql := "select id, email, password, created_at, updated_at from users where email = $1"

	var user models.User
	if err := ar.db.QueryRow(c, sql, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		return nil, err
	}
	return &user, nil
}

func (ar *AuthRepository) CreateAccount(c context.Context, user *models.User) error {
	tx, err := ar.db.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	qInsertIntoUser := "insert into users (email, password, created_at) values ($1, $2, now()) returning id"
	if err = tx.QueryRow(c, qInsertIntoUser, user.Email, user.Password).Scan(&user.ID); err != nil {
		return err
	}

	qInsertIntoProfile := "insert into profile (user_id, created_at) values ($1, now())"
	_, err = tx.Exec(c, qInsertIntoProfile, user.ID)
	if err != nil {
		return err
	}

	if err := tx.Commit(c); err != nil {
		return err
	}

	return nil
}
