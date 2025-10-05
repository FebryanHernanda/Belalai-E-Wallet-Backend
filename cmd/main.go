package main

import (
	"context"
	"log"

	_ "github.com/joho/godotenv/autoload"

	"github.com/Belalai-E-Wallet-Backend/internal/configs"
	"github.com/Belalai-E-Wallet-Backend/internal/routers"
)

// @title 											Belalai E-Wallet
// @version 										1.0
// @description 								E-wallet team belalai
// @host												127.0.0.1:3000/api/
// @securityDefinitions.apikey 	JWTtoken
// @in header
// @name Authorization
func main() {
	// Inisialization databae for this project
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

	// inisialization redish
	rdb := configs.InitRedis()
	cmd := rdb.Ping(context.Background())
	if cmd.Err() != nil {
		log.Println("failed ping on redis \nCause:", cmd.Err().Error())
		return
	}
	log.Println("Redis Connected")
	defer rdb.Close()

	// Inisialization engine gin, HTTP framework
	router := routers.InitRouter(db, rdb)
	router.Run(":2409")
}
