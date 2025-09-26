package routers

import (
	"github.com/Belalai-E-Wallet-Backend/internal/handler"
	"github.com/Belalai-E-Wallet-Backend/internal/middleware"
	"github.com/Belalai-E-Wallet-Backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitTransferRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	transferRouter := router.Group("/transfer")
	transferRepository := repository.NewTransferRepository(db, rdb)
	uh := handler.NewTransferHandler(transferRepository)

	transferRouter.GET("", middleware.VerifyToken(rdb), uh.FilterUser)
	transferRouter.POST("", middleware.VerifyToken(rdb), uh.TranferBalance)
}
