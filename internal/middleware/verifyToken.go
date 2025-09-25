package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Belalai-E-Wallet-Backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func VerifyToken(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// ambil token dari header
		bearerToken := ctx.GetHeader("Authorization")
		// Bearer token
		token := strings.Split(bearerToken, " ")[1]
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Silahkan login terlebih dahulu",
			})
			return
		}

		// !DO cek token from redish if it not blacklisted
		isBlacklisted, err := rdb.Get(ctx, "Belalai-E-wallet:blacklist:"+bearerToken).Result()
		if err == nil && isBlacklisted == "true" {
			log.Println("Token sudah logout, silahkan login kembali")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Token sudah logout, silahkan login kembali",
			})
			return
		} else if err != redis.Nil && err != nil {
			log.Println("Error when checking blacklist redis cache:", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Internal Server Error",
			})
			return
		}

		// verify token jwt
		var claims pkg.Claims
		if err := claims.VerifyToken(token); err != nil {
			if strings.Contains(err.Error(), jwt.ErrTokenInvalidIssuer.Error()) {
				log.Println("JWT Error.\nCause: ", err.Error())
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Silahkan login kembali",
				})
				return
			}
			if strings.Contains(err.Error(), jwt.ErrTokenExpired.Error()) {
				log.Println("JWT Error.\nCause: ", err.Error())
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Silahkan login kembali",
				})
				return
			}
			fmt.Println(jwt.ErrTokenExpired)
			log.Println("Internal Server Error.\nCause: ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Internal Server Error",
			})
			return
		}
		ctx.Set("claims", claims)
		ctx.Next()
	}
}
