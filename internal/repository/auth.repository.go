package repository

import (
	"context"
	"errors"
	"log"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/Belalai-E-Wallet-Backend/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type AuthRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewAuthRepository(db *pgxpool.Pool, rdb *redis.Client) *AuthRepository {
	return &AuthRepository{
		db:  db,
		rdb: rdb,
	}
}

func (ar *AuthRepository) GetEmail(c context.Context, email string) (*models.User, error) {
	sql := "select id, email, password, pin, created_at, updated_at from users where email = $1"

	var user models.User
	if err := ar.db.QueryRow(c, sql, email).Scan(&user.ID, &user.Email, &user.Password, &user.Pin, &user.CreatedAt, &user.UpdatedAt); err != nil {
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
		log.Println(err.Error())
		return err
	}

	qInsertIntoProfile := "insert into profile (user_id, created_at) values ($1, now())"
	_, err = tx.Exec(c, qInsertIntoProfile, user.ID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	qInsertIntoWallet := "insert into wallets (user_id, created_at) values ($1, now())"
	_, err = tx.Exec(c, qInsertIntoWallet, user.ID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if err := tx.Commit(c); err != nil {
		log.Println(err.Error())
		return err
	}

	// Handle caching after successful registration
	if err := utils.HandleRegistrationCache(c, *ar.rdb, int64(user.ID), user); err != nil {
		log.Println("Cache operation warning:", err)
		// Don't return error - registration was successful
	}

	return nil
}

// VerifyPassword: get hashed password by userId
func (ar *AuthRepository) VerifyPassword(c context.Context, userId int) (string, error) {
	var hashedPassword string
	sql := `SELECT password FROM users WHERE id = $1`

	err := ar.db.QueryRow(c, sql, userId).Scan(&hashedPassword)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return hashedPassword, nil
}

// UpdatePassword: update user password
func (ar *AuthRepository) UpdatePassword(c context.Context, userId int, hashedPassword string) error {
	sql := `UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`
	_, err := ar.db.Exec(c, sql, hashedPassword, userId)
	return err
}

// VerifyPIN: get hashed pin by userId
func (ar *AuthRepository) VerifyPIN(c context.Context, userId int) (string, error) {
	var hashedPIN string
	sql := `SELECT pin FROM users WHERE id = $1`

	err := ar.db.QueryRow(c, sql, userId).Scan(&hashedPIN)
	if err != nil {
		return "", err
	}

	return hashedPIN, nil
}

// UpdatePIN: update user pin
func (ar *AuthRepository) UpdatePIN(c context.Context, userId int, hashedPin string) error {
	sql := `UPDATE users SET pin = $1, updated_at = NOW() WHERE id = $2`
	_, err := ar.db.Exec(c, sql, hashedPin, userId)
	return err
}

// BlacklistToken: blacklist user token (logout)
func (ar *AuthRepository) BlacklistToken(c context.Context, token string) error {
	// use utils.BlackListTokenRedish for logout token
	if err := utils.BlackListTokenRedish(c, *ar.rdb, token); err != nil {
		log.Println("failed blacklist token, ", err)
		return err
	}
	// is success return nil
	return nil
}
