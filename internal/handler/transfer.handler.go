package handler

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/Belalai-E-Wallet-Backend/internal/repository"
	"github.com/Belalai-E-Wallet-Backend/internal/utils"
	"github.com/Belalai-E-Wallet-Backend/pkg"
	"github.com/gin-gonic/gin"
)

type TransferHandler struct {
	transRep *repository.TransferRepository
}

func NewTransferHandler(transRep *repository.TransferRepository) *TransferHandler {
	return &TransferHandler{transRep: transRep}
}

func (u *TransferHandler) FilterUser(ctx *gin.Context) {
	// default get all user if query is empty
	query := ctx.Query("search")
	// Make pagenation using query LIMIT dan OFFSET
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	limit := 10
	offset := (page - 1) * limit

	// use / call repository filter user
	users, err := u.transRep.FilterUser(ctx.Request.Context(), query, offset, limit, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: err.Error(),
		})
		return
	}

	// send data users as response
	ctx.JSON(http.StatusOK, models.ResponseData{
		Response: models.Response{
			IsSuccess: true,
			Code:      http.StatusOK,
		},
		Data: users,
	})
}

func (u *TransferHandler) TranferBalance(ctx *gin.Context) {
	// get user id from token
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

	// binding data JSON
	var body models.TransferBody
	if err := ctx.ShouldBind(&body); err != nil {
		log.Println("Failed binding data\nCause: ", err)
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      400,
			},
			Err: "Failed binding data...",
		})
		return
	}

	// get user hashed pin
	user, err := u.transRep.GetHashedPin(ctx.Request.Context(), userID)
	if err != nil {
		log.Println("failed get hashed pin \nCause: ", err)
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusBadRequest,
			},
			Err: "user sender is not found",
		})
		return
	}

	// compare the pin :
	// body.PinSender => from http body / input user
	// user.pin => from query GetHashedPin
	hc := pkg.NewHashConfig()
	isMatched, err := hc.CompareHashAndPassword(body.PinSender, user.Pin)
	if err != nil {
		log.Println("Internal Server Error.\nCause: ", err.Error())
		re := regexp.MustCompile("hash|crypto|argon2id|format")
		if re.Match([]byte(err.Error())) {
			log.Println("Error during Hashing")
		}
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: "internal server error",
		})
		return
	}

	// if pin is not match
	if !isMatched {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      400,
			},
			Err: "pin is incorrect",
		})
		return
	}

	// if match execute tranfer using func repo
	if err := u.transRep.TransferMoney(ctx.Request.Context(), userID, body); err != nil {
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: "internal server error",
		})
		return
	} else {
		ctx.JSON(http.StatusOK, models.Response{
			IsSuccess: true,
			Code:      http.StatusOK,
			Msg:       "transfer is success",
		})
	}
}
