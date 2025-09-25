package routers

import (
	"github.com/Belalai-E-Wallet-Backend/internal/handler"
	"github.com/Belalai-E-Wallet-Backend/internal/middleware"
	"github.com/Belalai-E-Wallet-Backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitEWalletRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	eWalletRouter := router.Group("/balance")
	eWalletRepository := repository.NewEWalletRepository(db)
	eWalletHandler := handler.NewEWalletHandler(eWalletRepository)

	eWalletRouter.GET("", middleware.VerifyToken(rdb), eWalletHandler.GetBalance)
}
