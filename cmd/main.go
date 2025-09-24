package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"

	"github.com/Belalai-E-Wallet-Backend/internal/configs"
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

	db, err := configs.InitDB()
	if err != nil {
		log.Println("FAILED TO CONNECT DB")
		return
	}

	defer db.Close()

	err = configs.PingDB(db)
	if err != nil {
		log.Println("PING TO DB FAILED", err.Error())
		return
	}

	log.Println("DB CONNECTED")

	// manual load ENV

	// Inisialization databae for this project

	// inisialization redish

	// Inisialization engine gin, HTTP framework
	// !NEED DB
	router := routers.InitRouter(db)
	router.Run(":2409")
}
