package handler

import (
	"log"
	"net/http"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/Belalai-E-Wallet-Backend/internal/repository"
	"github.com/Belalai-E-Wallet-Backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type ChartHandler struct {
	cr *repository.ChartRepository
}

func NewChartHandler(cr *repository.ChartRepository) *ChartHandler {
	return &ChartHandler{cr: cr}
}

// @Summary Mendapatkan data chart berdasarkan durasi filter
// @Description Mendapatkan data chart keuangan pengguna yang sudah diautentikasi (berdasarkan ID pengguna dari token JWT) dengan filter durasi tertentu.
// @Tags Chart
// @Accept json
// @Produce json
// @Param duration path string true "Durasi filter data chart (contoh: 'seven_days', 'five_weeks', 'twelve_months')"
// @Success 200 {object} models.ChartDataResponse "Data chart berhasil diambil"
// @Failure 401 {object} models.UnauthorizedResponse "Tidak terautentikasi (Unauthorized) - Token JWT tidak valid atau hilang"
// @Failure 500 {object} models.InternalErrorResponse "Kesalahan server internal"
// @Router /chart/{duration} [get]
// @Security JWTtoken
func (c *ChartHandler) GetDataChart(ctx *gin.Context) {
	// get filter type
	durationFilter := ctx.Param("duration")

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

	// get data from repository by user ID
	chartData, err := c.cr.GetChartData(ctx.Request.Context(), userID, durationFilter)
	if err != nil {
		log.Println("error cause: \n", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ChartDataResponse{
		Response: models.Response{
			IsSuccess: true,
			Code:      http.StatusOK,
		},
		Data: chartData,
	})
}
