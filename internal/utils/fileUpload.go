package utils

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"regexp"
	"time"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/gin-gonic/gin"
)

func FileUpload(c *gin.Context, file *multipart.FileHeader, prefix string) string {
	const maxSize = 2 * 1024 * 1024
	if file.Size > maxSize {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "File too large (max 2MB)",
		})
		return ""
	}

	ext := filepath.Ext(file.Filename)
	re := regexp.MustCompile(`(?i)\.(png|jpg|jpeg|webp)$`)
	if !re.MatchString(ext) {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "Invalid file type (only PNG, JPG, JPEG, WEBP allowed)",
		})
		return ""
	}

	filename := fmt.Sprintf("%s_%d%s", prefix, time.Now().UnixNano(), ext)
	location := filepath.Join("public", filename)

	if err := c.SaveUploadedFile(file, location); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			IsSuccess: false,
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to save file",
		})
		return ""
	}

	return filename
}
