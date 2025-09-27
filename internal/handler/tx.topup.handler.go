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

type TopUpHandler struct {
	topUpRepo *repository.TopUpRepository
}

func NewTopUpHandler(topUpRepo *repository.TopUpRepository) *TopUpHandler {
	return &TopUpHandler{topUpRepo: topUpRepo}
}

func (th *TopUpHandler) GetPaymentMethods(c *gin.Context) {
	methods, err := th.topUpRepo.FindAllPaymentMethods(c)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
				Msg:       "Failed to fetch payment methods",
			},
			Err: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ResponseData{
		Response: models.Response{
			IsSuccess: true,
			Code:      http.StatusOK,
			Msg:       "Get payment methods successfully",
		},
		Data: methods,
	})
}

func (th *TopUpHandler) CreateTopUp(c *gin.Context) {
	var req models.TopUp
	if err := c.ShouldBind(&req); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusBadRequest,
				Msg:       "Invalid request",
			},
			Err: err.Error(),
		})
		return
	}

	req.Status = models.TopUpPending

	newTopup, err := th.topUpRepo.CreateTopUp(c, &req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
				Msg:       "Failed to create top up",
			},
			Err: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.ResponseData{
		Response: models.Response{
			IsSuccess: true,
			Code:      http.StatusCreated,
			Msg:       "Create top up successfully",
		},
		Data: newTopup,
	})
}

func (th *TopUpHandler) MarkTopUpSuccess(c *gin.Context) {
	topupID, _ := strconv.Atoi(c.Param("id"))

	userID, err := utils.GetUserFromCtx(c)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusUnauthorized,
				Msg:       "Unauthorized",
			},
			Err: err.Error(),
		})
		return
	}

	walletID, err := th.topUpRepo.GetWalletIDByUserID(c, userID)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusNotFound,
				Msg:       "Wallet not found",
			},
			Err: err.Error(),
		})
		return
	}

	err = th.topUpRepo.UpdateStatusTopUp(c, topupID, models.TopUpSuccess)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
				Msg:       "Failed to update top up status",
			},
			Err: err.Error(),
		})
		return
	}

	topup, err := th.topUpRepo.GetTopUpByID(c, topupID)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
				Msg:       "Failed to get topup data",
			},
			Err: err.Error(),
		})
		return
	}

	err = th.topUpRepo.ApplyToWallet(c, walletID, topupID, topup.Amount)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
				Msg:       "Failed to update wallet balance",
			},
			Err: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Topup applied to wallet",
	})
}

func (th *TopUpHandler) CreateTopUpTransaction(c *gin.Context) {
	var req models.TopUpRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusBadRequest,
				Msg:       "Invalid request",
			},
			Err: err.Error(),
		})
		return
	}

	userID, err := utils.GetUserFromCtx(c)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusUnauthorized,
				Msg:       "Unauthorized",
			},
			Err: err.Error(),
		})
		return
	}

	topup := &models.TopUp{
		Amount:    req.Amount,
		Tax:       req.Tax,
		PaymentID: req.PaymentID,
	}

	newTopup, err := th.topUpRepo.CreateTopUpTransaction(c, topup, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
				Msg:       "Failed to process topup",
			},
			Err: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.ResponseData{
		Response: models.Response{
			IsSuccess: true,
			Code:      http.StatusCreated,
			Msg:       "Topup successful",
		},
		Data: newTopup,
	})
}
