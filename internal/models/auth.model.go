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
	Token      string `json:"token"`
	IsPinExist bool   `json:"is_pin_exist"`
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

type ForgotPasswordOrPINRequest struct {
	Email string `json:"email" form:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" form:"token" binding:"required"`
	NewPassword string `json:"new_password" form:"new_password" binding:"required"`
}

type ResetPINRequest struct {
	Token  string `json:"token" form:"token" binding:"required"`
	NewPIN string `json:"new_pin" form:"new_pin" binding:"required,min=6"`
}

type ResponseReset struct {
	Token string `json:"token"`
	Link  string `json:"link"`
}

type ConfirmPayment struct {
	PIN string `json:"pin" form:"pin" binding:"required"`
}
