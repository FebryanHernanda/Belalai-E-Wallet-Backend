package models

import "time"

type Profile struct {
	UserID         int        `db:"users_id"`
	Fullname       *string    `db:"fullname"`
	Phone          *string    `db:"phone"`
	ProfilePicture *string    `db:"profile_picture"`
	Email          *string    `db:"-"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
}

type ProfileRequest struct {
	Fullname       *string `json:"fullname" form:"fullname"`
	Phone          *string `json:"phone" form:"phone"`
	ProfilePicture *string `json:"profile_picture" form:"profile_picture"`
	Email          *string `json:"email" form:"email"`
}

type ProfileResponse struct {
	UserID         int        `json:"user_id"`
	Fullname       *string    `json:"fullname"`
	Phone          *string    `json:"phone"`
	ProfilePicture *string    `json:"profile_picture"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}
