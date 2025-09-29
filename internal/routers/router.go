package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	docs "github.com/Belalai-E-Wallet-Backend/docs"
	"github.com/Belalai-E-Wallet-Backend/internal/middleware"
	"github.com/Belalai-E-Wallet-Backend/internal/models"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(db *pgxpool.Pool, rdb *redis.Client) *gin.Engine {
	// inizialization engine gin
	router := gin.Default()
	router.Use(middleware.CORSMiddleware)

	// swaggo configuration
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// setup routing
	InitAuthRouter(router, db, rdb)

	InitTransferRouter(router, db, rdb)

	InitEWalletRouter(router, db, rdb)

	InitTransactionRouter(router, db, rdb)

	InitProfileRouter(router, db, rdb)

	InitTopUpRouter(router, db, rdb)

	InitChartRoouter(router, db, rdb)

	// make directori public accesible
	router.Static("/img", "public")

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, models.Response{
			IsSuccess: false,
			Code:      http.StatusNotFound,
			Msg:       "Page not found!",
		})
	})

	return router
}
