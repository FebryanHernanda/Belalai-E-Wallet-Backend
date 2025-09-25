package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/Belalai-E-Wallet-Backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type EWalletHandler struct {
	er *repository.EwalletRepository
}

func NewEWalletHandler(er *repository.EwalletRepository) *EWalletHandler {
	return &EWalletHandler{er: er}
}

// GetBalance
// @tags 				balance
// @router 	 		/balance 	[POST]
// @Summary 		Get user balance
// @Description Get balance for a specific user by user_id
// @Param 			body		body	 models.Balance  	true 		"Input user_id"
// @accept 			json
// @produce 		json
// @failure 		400			{object} 	models.BadRequestResponse "Bad Request"
// @failure 		500 		{object} 	models.InternalErrorResponse "Internal Server Error"
// @success 		200 		{object}  models.ResponseData "Success Response with Balance Data"
func (e *EWalletHandler) GetBalance(ctx *gin.Context) {
	var body models.Balance
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Println("error when binding \ncause", err)
		if strings.Contains(err.Error(), "required") {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      400,
				},
				Err: "user_id cannot be empty",
			})
			return
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: "Internal Server Error",
		})
		return
	}

	// Get balance from repository
	balance, err := e.er.GetBalance(ctx, body.User_id)
	if err != nil {
		log.Println("Error getting balance from repository.\nCause: ", err.Error())

		// Check if user not found
		if err.Error() == "user_id not found" {
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
