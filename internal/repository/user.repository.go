package repository

import (
	"context"
	"log"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// filter user by name or phone number
func (ur *UserRepository) FilterUser(c context.Context, query string, offset, limit int) ([]models.ProfileResponse, error) {
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
	// return users, and error nil if success
	return users, nil
}
