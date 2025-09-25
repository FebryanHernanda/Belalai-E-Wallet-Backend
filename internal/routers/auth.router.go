package routers

import (
	"github.com/Belalai-E-Wallet-Backend/internal/handler"
	"github.com/Belalai-E-Wallet-Backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitAuthRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authRouter := router.Group("/auth")
	authRepository := repository.NewAuthRepository(db, rdb)
	authHandler := handler.NewAuthHandler(authRepository)

	authRouter.POST("", authHandler.Login)
	authRouter.POST("/register", authHandler.Register)
	authRouter.DELETE("", authHandler.Logout)
	authRouter.PATCH("/update-pin", authHandler.UpdatePIN)
	authRouter.PATCH("/change-pin", authHandler.ChangePIN)
	authRouter.PATCH("/change-password", authHandler.ChangePassword)
}
