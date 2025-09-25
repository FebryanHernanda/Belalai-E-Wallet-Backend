package models

import "time"

type User struct {
	ID        int        `db:"id"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	Pin       *string    `db:"pin"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

type AuthRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type ChangePasswordRequest struct {
	OldPassword string `form:"old_password" json:"old_password" binding:"required"`
	NewPassword string `form:"new_password" json:"new_password" binding:"required"`
}

type ChangePINRequest struct {
	OldPIN string `form:"old_pin" json:"old_pin" binding:"required,min=6"`
	NewPIN string `form:"new_pin" json:"new_pin" binding:"required,min=6"`
}

type SetPINRequest struct {
	PIN string `form:"pin" json:"new_pin" binding:"required,min=6"`
}
