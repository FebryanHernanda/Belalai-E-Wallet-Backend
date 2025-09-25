package handler

import (
	"log"
	"net/http"
	"regexp"
	"strings"

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
			Token: jwtToken,
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
		ctx.JSON(http.StatusOK, models.Response{
			IsSuccess: true,
			Code:      200,
			Msg:       "user register successfully",
		})
	}
}

// ChangePassword
// @Tags        auth
// @Router      /auth/change-password [PUT]
// @Summary     Change current user password
// @Description Change the current user password by providing old password and new password
// @Accept      json
// @Produce     json
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
// @Router      /auth/change-pin [PUT]
// @Summary     Change current user PIN
// @Description Change the current user PIN by providing old PIN and new PIN (min 6 characters)
// @Accept      json
// @Produce     json
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
// @Router      /auth/update-pin [PUT]
// @Summary     Update current user PIN directly
// @Description Update the current user PIN directly without old PIN (min 6 characters)
// @Accept      json
// @Produce     json
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
