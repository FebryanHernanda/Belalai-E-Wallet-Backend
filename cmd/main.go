package main

import (
	"github.com/Belalai-E-Wallet-Backend/internal/routers"
)

// @title 											Belalai E-Wallet
// @version 										1.0
// @description 								E-wallet team belalai
// @host												127.0.0.1:2409
// @securityDefinitions.apikey 	JWTtoken
// @in header
// @name Authorization
func main() {
	// manual load ENV

	// Inisialization databae for this project

	// inisialization redish

	// Inisialization engine gin, HTTP framework
	// !NEED DB
	router := routers.InitRouter()
	router.Run(":2409")
}
