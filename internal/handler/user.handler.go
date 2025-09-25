package handler

import (
	"net/http"
	"strconv"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/Belalai-E-Wallet-Backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	urep *repository.UserRepository
}

func NewUserHandler(urep *repository.UserRepository) *UserHandler {
	return &UserHandler{urep: urep}
}

func (u *UserHandler) FilterUser(ctx *gin.Context) {
	// default get all user if query is empty
	query := ctx.Query("search")
	// Make pagenation using query LIMIT dan OFFSET
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	limit := 20
	offset := (page - 1) * limit

	// use / call repository filter user
	users, err := u.urep.FilterUser(ctx.Request.Context(), query, offset, limit)
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
			Page:      page,
		},
		Data: users,
	})
}
