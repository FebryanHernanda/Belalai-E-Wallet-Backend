// handler/transaction.go (Updated)
package handler

import (
	"log"
	"net/http"
	"strconv"

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
// @Description 	Get transaction history for authenticated user (excluding soft deleted)
// @accept 			json
// @produce 		json
// @Security 		BearerAuth
// @failure 		401			{object} 	models.UnauthorizedResponse "Unauthorized"
// @failure 		404			{object} 	models.NotFoundResponse "Transaction History Not Found"
// @failure 		500 		{object} 	models.InternalErrorResponse "Internal Server Error"
// @success 		200 		{object}  	models.ResponseData "Success Response with Transaction History Data"
func (th *TransactionHandler) GetTransactionHistory(ctx *gin.Context) {
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

	histories, err := th.tr.GetHistory(ctx, userID)
	if err != nil {
		log.Println("Error getting transaction history from repository.\nCause: ", err.Error())

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

// GetAllTransactionHistory - Get transfer + topup history
// @tags 			transaction
// @router 			/transaction/history/all 	[GET]
// @Summary 		Get all user transaction history (transfer + topup)
// @Description 	Get complete transaction history including transfers and topups for authenticated user
// @accept 			json
// @produce 		json
// @Security 		BearerAuth
// @failure 		401			{object} 	models.UnauthorizedResponse "Unauthorized"
// @failure 		404			{object} 	models.NotFoundResponse "Transaction History Not Found"
// @failure 		500 		{object} 	models.InternalErrorResponse "Internal Server Error"
// @success 		200 		{object}  	models.ResponseData "Success Response with Complete Transaction History Data"
func (th *TransactionHandler) GetAllTransactionHistory(ctx *gin.Context) {
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

	histories, err := th.tr.GetAllHistory(ctx, userID)
	if err != nil {
		log.Println("Error getting all transaction history from repository.\nCause: ", err.Error())

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
			Msg:       "Get All Transaction History Success",
		},
		Data: histories,
	})
}

// DeleteTransaction - Soft delete transaction
// @tags 			transaction
// @router 			/transaction/{id} 	[DELETE]
// @Summary 		Soft delete transaction
// @Description 	Soft delete transaction for authenticated user
// @accept 			json
// @produce 		json
// @Security 		BearerAuth
// @Param			id	path	int	true	"Transaction ID"
// @failure 		401			{object} 	models.UnauthorizedResponse "Unauthorized"
// @failure 		404			{object} 	models.NotFoundResponse "Transaction Not Found"
// @failure 		500 		{object} 	models.InternalErrorResponse "Internal Server Error"
// @success 		200 		{object}  	models.Response "Success Response"
func (th *TransactionHandler) DeleteTransaction(ctx *gin.Context) {
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

	// Get transaction ID from URL parameter
	transactionIDStr := ctx.Param("id")
	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      400,
			},
			Err: "Invalid transaction ID",
		})
		return
	}

	err = th.tr.SoftDeleteTransaction(ctx, transactionID, userID)
	if err != nil {
		log.Println("Error soft deleting transaction.\nCause: ", err.Error())

		if err.Error() == "transaction not found or user not authorized" {
			ctx.JSON(http.StatusNotFound, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      404,
				},
				Err: "Transaction not found or unauthorized",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: "Failed to delete transaction",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Transaction deleted successfully",
	})
}
