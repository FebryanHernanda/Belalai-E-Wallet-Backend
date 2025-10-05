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

// @Summary Memfilter daftar pengguna
// @Description Mendapatkan daftar pengguna dengan opsi pencarian dan paginasi. Digunakan untuk memilih pengguna tujuan transfer.
// @Tags Transfer
// @Accept json
// @Produce json
// @Param search query string false "Kata kunci pencarian nama atau nomor telepon"
// @Param page query int false "Nomor halaman untuk paginasi (default: 1)"
// @Success 200 {object} models.ResponseData{Data=models.ListprofileResponse} "Daftar pengguna berhasil diambil"
// @Failure 401 {object} models.UnauthorizedResponse "Tidak terautentikasi (Unauthorized) - Token JWT tidak valid atau hilang"
// @Failure 500 {object} models.InternalErrorResponse "Kesalahan server internal"
// @Router /transfer [get]
// @Security JWTtoken
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

// @Summary Melakukan transfer saldo
// @Description Melakukan proses transfer saldo dari pengguna yang terautentikasi ke pengguna tujuan, memerlukan verifikasi PIN.
// @Tags Transfer
// @Accept json
// @Produce json
// @Param request body models.TransferBody true "Detail transfer (ID penerima, jumlah, dan PIN pengirim)"
// @Success 200 {object} models.Response "Transfer berhasil"
// @Failure 400 {object} models.ErrorResponse "Permintaan tidak valid (contoh: data binding gagal, PIN salah, saldo tidak cukup, transfer ke diri sendiri)"
// @Failure 401 {object} models.UnauthorizedResponse "Tidak terautentikasi (Unauthorized) - Token JWT tidak valid atau hilang"
// @Failure 500 {object} models.InternalErrorResponse "Kesalahan server internal"
// @Router /transfer [post]
// @Security JWTtoken
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
		if err == repository.ErrNotEnoughBalance {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      http.StatusBadRequest,
				},
				Err: err.Error(),
			})
			return
		}
		if err == repository.ErrCantSendingToYourself {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      http.StatusBadRequest,
				},
				Err: err.Error(),
			})
			return
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
				Msg:       "internal server error",
			},
			Err: err.Error(),
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
