package handler

import (
	"log"
	"net/http"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/Belalai-E-Wallet-Backend/internal/repository"
	"github.com/Belalai-E-Wallet-Backend/internal/utils" // Import utils package
	"github.com/gin-gonic/gin"
)

type EWalletHandler struct {
	er *repository.EwalletRepository
}

func NewEWalletHandler(er *repository.EwalletRepository) *EWalletHandler {
	return &EWalletHandler{er: er}
}

// GetBalance
// @tags 			balance
// @router 	 		/balance 	[GET]
// @Summary 		Get user balance
// @Description 	Get balance for authenticated user
// @accept 			json
// @produce 		json
// @Security 		BearerAuth
// @failure 		401			{object} 	models.UnauthorizedResponse "Unauthorized"
// @failure 		404			{object} 	models.NotFoundResponse "User Not Found"
// @failure 		500 		{object} 	models.InternalErrorResponse "Internal Server Error"
// @success 		200 		{object}  models.ResponseData "Success Response with Balance Data"
func (e *EWalletHandler) GetBalance(ctx *gin.Context) {
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

	// Get balance from repository using user ID from token
	balance, err := e.er.GetBalance(ctx, userID)
	if err != nil {
		log.Println("Error getting balance from repository.\nCause: ", err.Error())

		// Check if user not found
		if err.Error() == "user id not found" {
			ctx.JSON(http.StatusNotFound, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      404,
				},
				Err: "User not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: "Failed to retrieve balance",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseData{
		Response: models.Response{
			IsSuccess: true,
			Code:      http.StatusOK,
			Msg:       "Get Balance Success",
		},
		Data: *balance,
	})
}
