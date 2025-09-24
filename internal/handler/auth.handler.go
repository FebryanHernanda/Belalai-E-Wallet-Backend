package handler

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/Belalai-E-Wallet-Backend/internal/repository"
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
	ctx.JSON(http.StatusOK, models.AuthResponse{
		Token: jwtToken,
	})
}
