package handler

import (
	"log"
	"net/http"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/Belalai-E-Wallet-Backend/internal/repository"
	"github.com/Belalai-E-Wallet-Backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileRepository *repository.ProfileRepository
}

func NewProfileHandler(pr *repository.ProfileRepository) *ProfileHandler {
	return &ProfileHandler{
		profileRepository: pr,
	}
}

// @Summary Mendapatkan detail profil pengguna
// @Description Mengambil detail profil pengguna yang sedang login (berdasarkan ID dari token JWT).
// @Tags Profile
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseData{Data=models.ProfileResponse} "Detail profil berhasil diambil"
// @Failure 401 {object} models.UnauthorizedResponse "Tidak terautentikasi (Unauthorized) - Token JWT tidak valid atau hilang"
// @Failure 404 {object} models.NotFoundResponse "Profil pengguna tidak ditemukan"
// @Router /profile [get]
// @Security JWTtoken
func (ph *ProfileHandler) GetProfile(c *gin.Context) {
	userId, err := utils.GetUserFromCtx(c)
	if err != nil {
		log.Println("error cause: ", err.Error())
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusUnauthorized,
			},
			Err: err.Error(),
		})
		return
	}

	profile, err := ph.profileRepository.GetProfile(c.Request.Context(), userId)
	if err != nil {
		log.Println("error cause: ", err.Error())
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusNotFound,
			},
			Err: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ResponseData{
		Response: models.Response{
			IsSuccess: true,
			Code:      http.StatusOK,
			Msg:       "Get profile successfully",
		},
		Data: models.ProfileResponse{
			UserID:         userId,
			Fullname:       profile.Fullname,
			Phone:          profile.Phone,
			ProfilePicture: profile.ProfilePicture,
			Email:          *profile.Email,
			CreatedAt:      &profile.CreatedAt,
			UpdatedAt:      profile.UpdatedAt,
		},
	})
}

// @Summary Memperbarui detail profil pengguna
// @Description Memperbarui nama lengkap, nomor telepon, email, dan/atau gambar profil pengguna yang sedang login.
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Param fullname formData string true "Nama lengkap pengguna"
// @Param phone formData string true "Nomor telepon pengguna"
// @Param email formData string false "Alamat email pengguna (opsional)"
// @Param profile_picture formData file false "Gambar profil baru (format file)"
// @Success 200 {object} models.Response "Profil berhasil diperbarui"
// @Failure 400 {object} models.ErrorResponse "Permintaan tidak valid (contoh: data form binding gagal, kesalahan upload file)"
// @Failure 401 {object} models.UnauthorizedResponse "Tidak terautentikasi (Unauthorized) - Token JWT tidak valid atau hilang"
// @Failure 500 {object} models.InternalErrorResponse "Kesalahan server internal"
// @Router /profile [patch]
// @Security JWTtoken
func (ph *ProfileHandler) UpdateProfile(c *gin.Context) {
	userId, err := utils.GetUserFromCtx(c)
	if err != nil {
		log.Println("error cause: ", err.Error())
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusUnauthorized,
			},
			Err: err.Error(),
		})
		return
	}

	var body models.ProfileRequest
	if err := c.ShouldBind(&body); err != nil {
		log.Println("error cause: ", err.Error())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusBadRequest,
			},
			Err: err.Error(),
		})
		return
	}

	var profilePic *string
	file, err := c.FormFile("profile_picture")
	if err == nil {
		if filename, err := utils.FileUpload(c, file, "avatar"); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      http.StatusBadRequest,
				},
				Err: err.Error(),
			})
			return
		} else {
			profilePic = &filename
		}
	}

	profile := models.Profile{
		UserID:         userId,
		Fullname:       body.Fullname,
		Phone:          body.Phone,
		ProfilePicture: profilePic,
		Email:          body.Email,
	}

	if err := ph.profileRepository.UpdateProfile(c.Request.Context(), &profile); err != nil {
		log.Println("error cause: ", err.Error())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Profile updated successfully",
	})
}

// @Summary Menghapus gambar profil
// @Description Menghapus gambar profil (avatar) pengguna yang sedang login.
// @Tags Profile
// @Accept json
// @Produce json
// @Success 200 {object} models.Response "Gambar profil berhasil dihapus"
// @Failure 401 {object} models.UnauthorizedResponse "Tidak terautentikasi (Unauthorized) - Token JWT tidak valid atau hilang"
// @Failure 500 {object} models.InternalErrorResponse "Kesalahan server internal saat menghapus avatar"
// @Router /profile/avatar [delete]
// @Security JWTtoken
func (ph *ProfileHandler) DeleteAvatar(c *gin.Context) {
	userId, err := utils.GetUserFromCtx(c)
	if err != nil {
		log.Println("error cause: ", err.Error())
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusUnauthorized,
			},
			Err: err.Error(),
		})
		return
	}

	if err := ph.profileRepository.DeleteAvatar(c.Request.Context(), userId); err != nil {
		log.Println("error cause: ", err.Error())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Deleted profile picture successfully",
	})
}
