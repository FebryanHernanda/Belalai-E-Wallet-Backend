package routers

import (
	"github.com/Belalai-E-Wallet-Backend/internal/handler"
	"github.com/Belalai-E-Wallet-Backend/internal/middleware"
	"github.com/Belalai-E-Wallet-Backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitTransactionRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	transactionRouter := router.Group("/transaction")
	transactionRepository := repository.NewTransactionRepository(db)
	transactionHandler := handler.NewTransactionHandler(transactionRepository)

	transactionRouter.GET("/history", middleware.VerifyToken(rdb), transactionHandler.GetTransactionHistory)
	transactionRouter.GET("/history/all", middleware.VerifyToken(rdb), transactionHandler.GetAllTransactionHistory)
	transactionRouter.DELETE("/:id", middleware.VerifyToken(rdb), transactionHandler.DeleteTransaction)
	transactionRouter.DELETE("/topup/:id", middleware.VerifyToken(rdb), transactionHandler.DeleteTopup)
}
