package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearer := ctx.GetHeader("Authorization")

		if bearer == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Message: "please login first",
			})
			return
		}
		res := strings.Split(bearer, " ")

		if len(res) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Message: "invalid bearer token",
			})
			return
		}

		if res[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Message: "invalid bearer token",
			})
		}

		var token pkg.JWTClaims
		if err := token.DecodeToken(res[1]); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
				Success: false,
				Message: "an server error occured",
			})
			return
		}

		blacklistKey := os.Getenv("REDIS_USER") + ":blacklist:" + token.ID
		exists, err := rdb.Exists(ctx, blacklistKey).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			log.Println(err)
			ctx.AbortWithStatusJSON(500, dto.Response{
				Success: false,
				Message: "Internal server error while checking token validity",
				Data:    nil,
			})
			return
		}

		if exists > 0 {
			ctx.AbortWithStatusJSON(401, dto.Response{
				Success: false,
				Message: "Token has been revoked or logged out",
				Data:    nil,
			})
			return
		}

		ctx.Set("token", token)
		ctx.Next()
	}
}
