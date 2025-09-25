package routers

import (
	"github.com/Belalai-E-Wallet-Backend/internal/handler"
	"github.com/Belalai-E-Wallet-Backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitProfileRouter(router *gin.Engine, db *pgxpool.Pool) {
	profile := router.Group("/profile")
	profileRepo := repository.NewProfileRepository(db)
	profileHandler := handler.NewProfileHandler(profileRepo)

	profile.GET("", profileHandler.GetProfile)
	profile.PATCH("", profileHandler.UpdateProfile)
	profile.DELETE("/avatar", profileHandler.DeleteAvatar)
}
