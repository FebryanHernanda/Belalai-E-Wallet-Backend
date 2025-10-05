package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

func FileUpload(c *gin.Context, file *multipart.FileHeader, prefix string) (string, error) {
	const maxSize = 2 * 1024 * 1024
	if file.Size > maxSize {
		return "", errors.New("file too large (max 2MB)")
	}

	ext := filepath.Ext(file.Filename)
	re := regexp.MustCompile(`(?i)\.(png|jpg|jpeg|webp)$`)
	if !re.MatchString(ext) {
		return "", errors.New("invalid file type (only PNG, JPG, JPEG, WEBP allowed)")
	}

	filename := fmt.Sprintf("%s_%d%s", prefix, time.Now().UnixNano(), ext)
	location := filepath.Join("public", filename)

	if err := c.SaveUploadedFile(file, location); err != nil {
		return "", errors.New("failed to save file")
	}

	return filename, nil
}
