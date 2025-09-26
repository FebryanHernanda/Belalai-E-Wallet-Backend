package models

import (
	"mime/multipart"
	"time"
)

type Profile struct {
	UserID         int        `db:"user_id"`
	Fullname       *string    `db:"fullname"`
	Phone          *string    `db:"phone"`
	ProfilePicture *string    `db:"profile_picture"`
	Email          *string    `db:"email"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
}

type ProfileRequest struct {
	Fullname       *string               `json:"fullname" form:"fullname"`
	Phone          *string               `json:"phone" form:"phone"`
	ProfilePicture *multipart.FileHeader `form:"profile_picture"`
	Email          *string               `json:"email" form:"email"`
}

type ProfileResponse struct {
	UserID         int        `json:"user_id"`
	Fullname       *string    `json:"fullname"`
	Phone          *string    `json:"phone"`
	ProfilePicture *string    `json:"profile_picture"`
	Email          string     `json:"email"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}
