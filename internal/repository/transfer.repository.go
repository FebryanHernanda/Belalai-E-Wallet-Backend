package repository

import (
	"context"
	"log"
	"time"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/Belalai-E-Wallet-Backend/internal/utils"
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
func (ur *TransferRepository) FilterUser(c context.Context, query string, offset, limit int) ([]models.ProfileResponse, error) {

	// Get cached filter user, before accesing database
	// Get and renew cache only for page 1 only (offset 0)
	rdbKey := "Belalai-E-wallet:filter-user"
	if offset == 0 && query == "" {
		cachedUser, err := utils.RedisGetData[[]models.ProfileResponse](c, *ur.rdb, rdbKey)
		if err != nil {
			log.Println("Redis error :", err)
		} else if cachedUser != nil && len(*cachedUser) > 0 {
			return *cachedUser, nil
		}
	}

	// if there is no key/ no cached data
	// get filter user from database
	sql := `SELECT user_id, profile_picture, fullname, phone FROM profile
    WHERE fullname ILIKE $1 OR phone ILIKE $1
		LIMIT $2 OFFSET $3`

	searchQuery := "%" + query + "%"
	values := []any{searchQuery, limit, offset}

	// query user
	rows, err := ur.db.Query(c, sql, values...)
	if err != nil {
		log.Println("internal server error : ", err.Error())
		return []models.ProfileResponse{}, err
	}
	defer rows.Close()

	// processing data / read rows
	var users []models.ProfileResponse
	for rows.Next() {
		var user models.ProfileResponse
		if err := rows.Scan(&user.UserID, &user.ProfilePicture, &user.Fullname, &user.Phone); err != nil {
			log.Println("Scan Error, ", err.Error())
			return []models.ProfileResponse{}, err
		}
		users = append(users, user)
	}

	// make cache filter user after query data from database
	// Get and renew cache only for page 1 only (offset 0)
	if offset == 0 && query == "" {
		if err := utils.RedisRenewData(c, *ur.rdb, rdbKey, users, 10*time.Minute); err != nil {
			log.Println("Failed to renew Redis cache:", err.Error())
		}
	}

	// return users, and error nil if success
	return users, nil
}
