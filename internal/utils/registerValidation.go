package utils

import (
	"errors"
	"log"
	"regexp"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
)

func RegisterValidation(body models.AuthRequest) error {
	// cek format email
	// ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
	regexEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !regexEmail.Match([]byte(body.Email)) {
		return errors.New("email format is wrong")
	}

	// cek format password
	// harus : huruf, angka, simbol, 8 karakter
	log.Println(body.Password)
	islengEight := len(body.Password) >= 8
	isNotHvSymbl := regexp.MustCompile(`[!@#$%^&*/><]`).MatchString(body.Password)
	isNotHvChar := regexp.MustCompile(`[a-zA-Z]`).MatchString(body.Password)
	isNotHvDigit := regexp.MustCompile(`\d`).MatchString(body.Password)

	log.Println(isNotHvSymbl, isNotHvChar, isNotHvDigit, islengEight)
	if !isNotHvChar || !isNotHvSymbl || !isNotHvDigit || !islengEight {
		return errors.New("password must contain : character, digit, symbol, 8 characters leng")
	}
	return nil
}
