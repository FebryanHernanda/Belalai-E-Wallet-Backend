package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}

func (pr *ProfileRepository) GetProfile(c context.Context, userId int) (*models.Profile, error) {
	sql := `
		select user_id, profile_picture, fullname, phone, created_at, updated_at from profile where user_id = $1
	`

	var p models.Profile
	if err := pr.db.QueryRow(c, sql, userId).Scan(&p.UserID, &p.ProfilePicture, &p.Fullname, &p.Phone, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}

	return &p, nil
}

func (pr *ProfileRepository) UpdateProfile(c context.Context, profile *models.Profile) error {
	tx, err := pr.db.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	if profile.Email != nil {
		qSetEmail := "UPDATE users SET email = $1, updated_at = now() WHERE id = $2"
		if _, err := tx.Exec(c, qSetEmail, *profile.Email, profile.UserID); err != nil {
			return err
		}
	}

	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if profile.Fullname != nil {
		setClauses = append(setClauses, fmt.Sprintf("fullname = $%d", argPos))
		args = append(args, *profile.Fullname)
		argPos++
	}
	if profile.Phone != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", argPos))
		args = append(args, *profile.Phone)
		argPos++
	}
	if profile.ProfilePicture != nil {
		setClauses = append(setClauses, fmt.Sprintf("profile_picture = $%d", argPos))
		args = append(args, *profile.ProfilePicture)
		argPos++
	}

	if len(setClauses) > 0 {
		setClauses = append(setClauses, "updated_at = now()")

		query := "UPDATE profile SET " + strings.Join(setClauses, ", ")
		query += fmt.Sprintf(" WHERE user_id = $%d", argPos)
		args = append(args, profile.UserID)

		if _, err := tx.Exec(c, query, args...); err != nil {
			return err
		}
	}

	if err := tx.Commit(c); err != nil {
		return err
	}

	return nil
}

func (pr *ProfileRepository) DeleteAvatar(c context.Context, userId int) error {
	qDeleteAvatar := "update profile set profile_picture = null, updated_at = now() where user_id = $1"
	_, err := pr.db.Exec(c, qDeleteAvatar, userId)
	return err
}
