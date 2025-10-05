package routers

import (
	"github.com/Belalai-E-Wallet-Backend/internal/handler"
	"github.com/Belalai-E-Wallet-Backend/internal/middleware"
	"github.com/Belalai-E-Wallet-Backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitTopUpRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	topupRouter := router.Group("/topup")
	topupRepository := repository.NewTopUpRepository(db)
	topupHandler := handler.NewTopUpHandler(topupRepository)

	topupRouter.GET("/methods", middleware.VerifyToken(rdb), topupHandler.GetPaymentMethods)
	// topupRouter.POST("", middleware.VerifyToken(rdb), topupHandler.CreateTopUp)
	// topupRouter.PATCH("/:id/success", middleware.VerifyToken(rdb), topupHandler.MarkTopUpSuccess)

	// {
	// 	"amount": 100000,
	// 	"tax": 2500,
	// 	"payment_id": 2
	// }
	topupRouter.POST("", middleware.VerifyToken(rdb), topupHandler.CreateTopUpTransaction)
}
