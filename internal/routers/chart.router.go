package routers

import (
	"github.com/Belalai-E-Wallet-Backend/internal/handler"
	"github.com/Belalai-E-Wallet-Backend/internal/middleware"
	"github.com/Belalai-E-Wallet-Backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitChartRoouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	chartRouter := router.Group("/chart")
	chartRepository := repository.NewChartRepository(db)
	chartHandler := handler.NewChartHandler(chartRepository)

	chartRouter.GET("/:duration", middleware.VerifyToken(rdb), chartHandler.GetDataChart)
}
