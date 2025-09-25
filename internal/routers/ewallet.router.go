package routers

import (
	"github.com/Belalai-E-Wallet-Backend/internal/handler"
	"github.com/Belalai-E-Wallet-Backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitEWalletRouter(router *gin.Engine, db *pgxpool.Pool) {
	eWalletRouter := router.Group("/balance")
	eWalletRepository := repository.NewEWalletRepository(db)
	eWalletHandler := handler.NewEWalletHandler(eWalletRepository)

	eWalletRouter.GET("", eWalletHandler.GetBalance)
}
