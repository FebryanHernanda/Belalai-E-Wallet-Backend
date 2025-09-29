package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/Belalai-E-Wallet-Backend/internal/repository"
	"github.com/Belalai-E-Wallet-Backend/internal/utils"
	"github.com/Belalai-E-Wallet-Backend/pkg"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	ar *repository.AuthRepository
}

func NewAuthHandler(ar *repository.AuthRepository) *AuthHandler {
	return &AuthHandler{ar: ar}
}

// Login
// @tags 				login
// @router 	 		/auth 	[POST]
// @Summary 		Login registered user
// @Description login using email and password and return as response with JWT token
// @Param 			body		body	 models.AuthRequest  	true 		"Input email and password"
// @accept 			json
// @produce 		json
// @failure 		400			{object} 	models.BadRequestResponse "Bad Request"
// @failure 		500 		{object} 	models.InternalErrorResponse "Internal Server Error"
// @success 		200 		{object}  models.AuthResponse
func (a *AuthHandler) Login(ctx *gin.Context) {
	var body models.AuthRequest
	if err := ctx.ShouldBind(&body); err != nil {
		// check if failed binding bcs input not match with model require
		log.Println("error when binding \ncause", err)
		if strings.Contains(err.Error(), "required") {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      400,
				},
				Err: "Email or password cannot be empty",
			})
			return
		}
		// else binding error because server, log the error
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: "internal server error",
		})
		return
	}
	// get userdata and validate user
	user, err := a.ar.GetEmail(ctx.Request.Context(), body.Email)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      400,
				},
				Err: "Email or Password is incorrect",
			})
			return
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
	// check if pin is exist
	var isPinExist bool
	if user.Pin == nil {
		isPinExist = false
	} else {
		isPinExist = true
	}

	// compare the password :
	// body.password => from http body / input user
	// user.Password => from query GetUserWithEmail
	hc := pkg.NewHashConfig()
	isMatched, err := hc.CompareHashAndPassword(body.Password, user.Password)
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
	// if not match sen https status as response
	if !isMatched {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      400,
			},
			Err: "Email or Password is incorrect",
		})
		return
	}
	// If match, generate jwt token and send as response
	claim := pkg.NewJWTClaims(user.ID, "user")
	jwtToken, err := claim.GenToken()
	if err != nil {
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: "internal server error",
		})
		return
	}
	// return token as response success
	ctx.JSON(http.StatusOK, models.ResponseData{
		Response: models.Response{
			IsSuccess: true,
			Code:      http.StatusOK,
			Msg:       "login successfully",
		},
		Data: models.AuthResponse{
			Token:      jwtToken,
			IsPinExist: isPinExist,
		},
	})
}

// Register
// @Tags				/auth/register [POST]
// @Summary 		Register new user
// @Description	Register new user with input email & password
// @Param				body		body 		 models.AuthRequest 	true		"Input email and password new user"
// @accept			json
// @produce			json
// @failure 		400			{object} 	models.BadRequestResponse "Bad Request"
// @failure 		500 		{object} 	models.InternalErrorResponse "Internal Server Error"
// @success			200			{object}  models.Response
func (a *AuthHandler) Register(ctx *gin.Context) {
	var body models.AuthRequest

	// Binding data and show if there is error when binding data
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: "Failed binding data ...",
		})
		return
	}

	// validate user register
	if err := utils.RegisterValidation(body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      400,
			},
			Err: err.Error(),
		})
		return
	} else {
		// hash new password
		hc := pkg.NewHashConfig()
		hc.UseRecommended()
		hash, err := hc.GenHash(body.Password)
		if err != nil {
			log.Println("Failed hash new password ...", err)
			ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      500,
				},
				Err: err.Error(),
			})
			return
		}
		// if inputs is already valid format,
		// input and check if the email already registered
		user := models.User{
			Email:    body.Email,
			Password: hash,
		}
		if err := a.ar.CreateAccount(ctx.Request.Context(), &user); err != nil {
			log.Println("error cause: ", err)
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      400,
				},
				Err: err.Error(),
			})
			return
		}

		go func() {
			err := utils.Send(utils.SendOptions{
				To:         []string{body.Email},
				Subject:    "Welcome to E-Wallet!",
				Body:       fmt.Sprintf("<h2>Hello %s!</h2><p>Terima kasih sudah mendaftar di E-Wallet.</p>", body.Email),
				BodyIsHTML: true,
			})
			if err != nil {
				log.Println("Failed to send registration email:", err)
			} else {
				log.Printf("Email registration sent to %s\n", body.Email)
			}
		}()

		ctx.JSON(http.StatusOK, models.Response{
			IsSuccess: true,
			Code:      200,
			Msg:       "user register successfully",
		})
	}
}

// ChangePassword
// @Tags        auth
// @Router      /auth/change-password [PATCH]
// @Summary     Change current user password
// @Description Change the current user password by providing old password and new password
// @Accept      json
// @Produce     json
// @Security    JWTtoken
// @Param       body  body      models.ChangePasswordRequest true "Old and New Password"
// @Success     200   {object}  models.Response
// @Failure     400   {object}  models.BadRequestResponse "Bad Request"
// @Failure     401   {object}  models.ErrorResponse "Unauthorized"
// @Failure     500   {object}  models.InternalErrorResponse "Internal Server Error"
func (a *AuthHandler) ChangePassword(ctx *gin.Context) {
	userId, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		log.Println("error cause: ", err)
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusUnauthorized,
			},
			Err: err.Error(),
		})
		return
	}

	var body models.ChangePasswordRequest
	if err := ctx.ShouldBind(&body); err != nil {
		log.Println("error cause: ", err)
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusBadRequest,
			},
			Err: err.Error(),
		})
		return
	}

	pwdFromDB, err := a.ar.VerifyPassword(ctx, userId)
	if err != nil {
		log.Println("error cause: ", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: err.Error(),
		})
		return
	}

	hashConfig := pkg.NewHashConfig()
	ok, _ := hashConfig.CompareHashAndPassword(body.OldPassword, pwdFromDB)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusUnauthorized,
			},
			Err: "invalid old password",
		})
		return
	}

	hashConfig.UseRecommended()
	hashedPwd, err := hashConfig.GenHash(body.NewPassword)
	if err != nil {
		log.Println("error cause: ", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: err.Error(),
		})
		return
	}

	if err := a.ar.UpdatePassword(ctx, userId, hashedPwd); err != nil {
		log.Println("error cause: ", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Password changed successfully",
	})
}

// ChangePIN
// @Tags        auth
// @Router      /auth/change-pin [PATCH]
// @Summary     Change current user PIN
// @Description Change the current user PIN by providing old PIN and new PIN (min 6 characters)
// @Accept      json
// @Produce     json
// @Security    JWTtoken
// @Param       body  body      models.ChangePINRequest true "Old and New PIN"
// @Success     200   {object}  models.Response
// @Failure     400   {object}  models.BadRequestResponse "Bad Request"
// @Failure     401   {object}  models.ErrorResponse "Unauthorized"
// @Failure     500   {object}  models.InternalErrorResponse "Internal Server Error"
func (a *AuthHandler) ChangePIN(ctx *gin.Context) {
	userId, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		log.Println("error cause: ", err)
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusUnauthorized,
			},
			Err: err.Error(),
		})
		return
	}

	var body models.ChangePINRequest
	if err := ctx.ShouldBind(&body); err != nil {
		log.Println("error cause: ", err)
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusBadRequest,
			},
			Err: err.Error(),
		})
		return
	}

	pinDB, err := a.ar.VerifyPIN(ctx, userId)
	if err != nil {
		log.Println("error cause: ", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: err.Error(),
		})
		return
	}

	hashConfig := pkg.NewHashConfig()
	ok, _ := hashConfig.CompareHashAndPassword(body.OldPIN, pinDB)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusUnauthorized,
			},
			Err: "invalid old pin",
		})
		return
	}

	hashConfig.UseRecommended()
	hashedPin, err := hashConfig.GenHash(body.NewPIN)
	if err != nil {
		log.Println("error cause: ", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: err.Error(),
		})
		return
	}

	if err := a.ar.UpdatePIN(ctx, userId, hashedPin); err != nil {
		log.Println("error cause: ", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "PIN changed successfully",
	})
}

// UpdatePIN
// @Tags        auth
// @Router      /auth/update-pin [PATCH]
// @Summary     Update current user PIN directly
// @Description Update the current user PIN directly without old PIN (min 6 characters)
// @Accept      json
// @Produce     json
// @Security    JWTtoken
// @Param       body  body      models.SetPINRequest true "New PIN"
// @Success     200   {object}  models.Response
// @Failure     400   {object}  models.BadRequestResponse "Bad Request"
// @Failure     401   {object}  models.ErrorResponse "Unauthorized"
// @Failure     500   {object}  models.InternalErrorResponse "Internal Server Error"
func (a *AuthHandler) UpdatePIN(ctx *gin.Context) {
	userId, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		log.Println("error cause: ", err)
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusUnauthorized,
			},
			Err: err.Error(),
		})
		return
	}

	var body models.SetPINRequest
	if err := ctx.ShouldBind(&body); err != nil {
		log.Println("error cause: ", err)
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusBadRequest,
			},
			Err: err.Error(),
		})
		return
	}

	hashConfig := pkg.NewHashConfig()
	hashConfig.UseRecommended()

	hashedPin, err := hashConfig.GenHash(body.PIN)
	if err != nil {
		log.Println("error cause: ", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: err.Error(),
		})
		return
	}

	if err := a.ar.UpdatePIN(ctx, userId, hashedPin); err != nil {
		log.Println("error cause: ", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "PIN updated successfully",
	})
}

// Logout
// @Tags			logout
// @Router			/auth [DELETE]
// @Summary 		Logout user by blacklist their token
// @Description	Logout user by blacklist their token on redis
// @Security 		JWTtoken
// @produce			json
// @failure 		500 	{object} 	models.InternalErrorResponse "Internal Server Error"
// @success			200 	{object}	models.Response
func (a *AuthHandler) Logout(ctx *gin.Context) {
	// get token user for logout
	bearerToken := ctx.GetHeader("Authorization")

	if err := a.ar.BlacklistToken(ctx.Request.Context(), bearerToken); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: err.Error(),
		})
		return
	} else {
		ctx.JSON(http.StatusOK, models.Response{
			IsSuccess: true,
			Code:      200,
			Msg:       "Logout successfully",
		})
	}
}

// ForgotPassword
// @Tags        auth
// @Summary     Request password reset
// @Description Send a reset password link with token to the user's email
// @Accept      json
// @Produce     json
// @Param       body body models.ForgotPasswordOrPINRequest true "User email"
// @Success     200 {object} models.Response
// @Failure     400 {object} models.ErrorResponse "Invalid email format"
// @Failure     500 {object} models.InternalErrorResponse "Internal Server Error"
// @Router      /auth/forgot-password [post]
func (a *AuthHandler) ForgotPassword(ctx *gin.Context) {
	var body models.ForgotPasswordOrPINRequest
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusBadRequest,
			},
			Err: "Invalid email format",
		})
		return
	}

	user, err := a.ar.GetEmailForSMPT(ctx.Request.Context(), body.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: "User not registered",
		})
		return
	}

	token, err := utils.GenerateRandomToken(32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: "Failed to generate token",
		})
		return
	}

	key := "reset:pwd:" + token
	if err := a.ar.SaveResetToken(ctx, key, fmt.Sprintf("%d", user.ID), 15*time.Minute); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: "Failed to save reset token",
		})
		return
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)
	go utils.Send(utils.SendOptions{
		To:         []string{body.Email},
		Subject:    "Reset Password",
		Body:       fmt.Sprintf("<p>Klik link berikut untuk reset password Anda:</p><p><a href='%s'>Reset Password</a></p>", resetLink),
		BodyIsHTML: true,
	})

	ctx.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Link reset password was sent to email",
	})
}

// ResetPassword
// @Tags        auth
// @Summary     Reset user password
// @Description Reset user password using the token received via email
// @Accept      json
// @Produce     json
// @Param       body body models.ResetPasswordRequest true "Token and new password"
// @Success     200 {object} models.Response
// @Failure     400 {object} models.ErrorResponse "Invalid or expired token"
// @Failure     500 {object} models.InternalErrorResponse "Internal Server Error"
// @Router      /auth/reset-password [post]
func (a *AuthHandler) ResetPassword(ctx *gin.Context) {
	var body models.ResetPasswordRequest
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusBadRequest,
			},
			Err: "token and new password required",
		})
		return
	}

	key := "reset:pwd:" + body.Token
	userIdStr, err := a.ar.GetResetToken(ctx, key)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusBadRequest,
			},
			Err: "invalid or expired token",
		})
		return
	}
	defer a.ar.DeleteResetToken(ctx, key)

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: "failed to parse user id",
		})
		return
	}

	// hash password baru
	hc := pkg.NewHashConfig()
	hc.UseRecommended()
	hashedPwd, err := hc.GenHash(body.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: "failed to hash password",
		})
		return
	}

	if err := a.ar.UpdatePassword(ctx, userId, hashedPwd); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: "failed to update password",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Password reset successfully",
	})
}

// ForgotPIN
// @Tags        auth
// @Summary     Request PIN reset
// @Description Send a reset PIN link with token to the user's email
// @Accept      json
// @Produce     json
// @Param       body body models.ForgotPasswordOrPINRequest true "User email"
// @Success     200 {object} models.Response
// @Failure     400 {object} models.ErrorResponse "Invalid email format"
// @Failure     500 {object} models.InternalErrorResponse "Internal Server Error"
// @Router      /auth/forgot-pin [post]
func (a *AuthHandler) ForgotPIN(ctx *gin.Context) {
	var body models.ForgotPasswordOrPINRequest
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusBadRequest,
			},
			Err: "invalid email format",
		})
		return
	}

	user, err := a.ar.GetEmailForSMPT(ctx.Request.Context(), body.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: "User not registered",
		})
		return
	}

	token, err := utils.GenerateRandomToken(32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: "failed to generate token",
		})
		return
	}

	key := "reset:pin:" + token
	if err := a.ar.SaveResetToken(ctx, key, fmt.Sprintf("%d", user.ID), 15*time.Minute); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: "failed to save reset token",
		})
		return
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	resetLink := fmt.Sprintf("%s/reset-pin?token=%s", frontendURL, token)
	go utils.Send(utils.SendOptions{
		To:         []string{body.Email},
		Subject:    "Reset PIN",
		Body:       fmt.Sprintf("<p>Klik link berikut untuk reset PIN Anda:</p><p><a href='%s'>Reset PIN</a></p>", resetLink),
		BodyIsHTML: true,
	})

	ctx.JSON(http.StatusOK, models.ResponseData{
		Response: models.Response{
			IsSuccess: true,
			Code:      http.StatusOK,
			Msg:       "Link reset password was sent to email",
		},
	})
}

// ResetPIN
// @Tags        auth
// @Summary     Reset user PIN
// @Description Reset user PIN using the token received via email
// @Accept      json
// @Produce     json
// @Param       body body models.ResetPINRequest true "Token and new_pin"
// @Success     200 {object} models.Response
// @Failure     400 {object} models.ErrorResponse "Invalid or expired token"
// @Failure     500 {object} models.InternalErrorResponse "Internal Server Error"
// @Router      /auth/reset-pin [post]
func (a *AuthHandler) ResetPIN(ctx *gin.Context) {
	var body models.ResetPINRequest
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusBadRequest,
			},
			Err: "token and new pin required",
		})
		return
	}

	key := "reset:pin:" + body.Token
	userIdStr, err := a.ar.GetResetToken(ctx, key)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusBadRequest,
			},
			Err: "invalid or expired token",
		})
		return
	}
	defer a.ar.DeleteResetToken(ctx, key)

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: "failed to parse user id",
		})
		return
	}

	hc := pkg.NewHashConfig()
	hc.UseRecommended()
	hashedPin, err := hc.GenHash(body.NewPIN)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: "failed to hash pin",
		})
		return
	}

	if err := a.ar.UpdatePIN(ctx, userId, hashedPin); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: "failed to update pin",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "PIN reset successfully",
	})
}

// ConfirmPIN
// @Tags        auth
// @Summary     Confirm user PIN
// @Description Verify user's PIN before processing a payment
// @Security    JWTtoken
// @Accept      json
// @Produce     json
// @Param       body body models.ConfirmPayment true "PIN confirmation"
// @Success     200 {object} models.Response
// @Failure     400 {object} models.ErrorResponse "Invalid request"
// @Failure     401 {object} models.ErrorResponse "Invalid PIN or unauthorized"
// @Failure     500 {object} models.InternalErrorResponse "Internal Server Error"
// @Router      /auth/confirm-pin [post]
func (a *AuthHandler) ConfirmPIN(ctx *gin.Context) {
	userId, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusUnauthorized,
			},
			Err: err.Error(),
		})
		return
	}

	var body models.ConfirmPayment
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusBadRequest,
			},
			Err: err.Error(),
		})
		return
	}

	pinDB, err := a.ar.VerifyPIN(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
			},
			Err: err.Error(),
		})
		return
	}

	hashConfig := pkg.NewHashConfig()
	ok, err := hashConfig.CompareHashAndPassword(body.PIN, pinDB)
	if err != nil || !ok {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      http.StatusUnauthorized,
			},
			Err: "invalid pin",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "PIN verified successfully, payment confirmed",
	})
}
