package models

type Balance struct {
	User_id int `db:"user_id"`
	Balance int `db:"balance"`
}
