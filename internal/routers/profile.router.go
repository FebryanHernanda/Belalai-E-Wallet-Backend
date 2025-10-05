package routers

import (
	"github.com/Belalai-E-Wallet-Backend/internal/handler"
	"github.com/Belalai-E-Wallet-Backend/internal/middleware"
	"github.com/Belalai-E-Wallet-Backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitProfileRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	profile := router.Group("/profile")
	profileRepo := repository.NewProfileRepository(db, *rdb)
	profileHandler := handler.NewProfileHandler(profileRepo)

	profile.GET("", middleware.VerifyToken(rdb), profileHandler.GetProfile)
	profile.PATCH("", middleware.VerifyToken(rdb), profileHandler.UpdateProfile)
	profile.DELETE("/avatar", middleware.VerifyToken(rdb), profileHandler.DeleteAvatar)
}
