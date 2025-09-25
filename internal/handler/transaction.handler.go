// handler/transaction.go
package handler

import (
	"log"
	"net/http"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/Belalai-E-Wallet-Backend/internal/repository"
	"github.com/Belalai-E-Wallet-Backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	tr *repository.TransactionRepository
}

func NewTransactionHandler(tr *repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{tr: tr}
}

// GetTransactionHistory
// @tags 			transaction
// @router 			/transaction/history 	[GET]
// @Summary 		Get user transaction history
// @Description 	Get transaction history for authenticated user
// @accept 			json
// @produce 		json
// @Security 		BearerAuth
// @failure 		401			{object} 	models.UnauthorizedResponse "Unauthorized"
// @failure 		404			{object} 	models.NotFoundResponse "Transaction History Not Found"
// @failure 		500 		{object} 	models.InternalErrorResponse "Internal Server Error"
// @success 		200 		{object}  	models.ResponseData "Success Response with Transaction History Data"
func (th *TransactionHandler) GetTransactionHistory(ctx *gin.Context) {
	// Get user ID from JWT token in context
	userID, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		log.Println("Error getting user from context.\nCause: ", err.Error())
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      401,
			},
			Err: "Unauthorized: " + err.Error(),
		})
		return
	}

	// Get transaction history from repository using user ID from token
	histories, err := th.tr.GetHistory(ctx, userID)
	if err != nil {
		log.Println("Error getting transaction history from repository.\nCause: ", err.Error())

		// Check if no transactions found
		if err.Error() == "no transactions found" {
			ctx.JSON(http.StatusNotFound, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      404,
				},
				Err: "Transaction history not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: "Failed to retrieve transaction history",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseData{
		Response: models.Response{
			IsSuccess: true,
			Code:      http.StatusOK,
			Msg:       "Get Transaction History Success",
		},
		Data: histories,
	})
}
