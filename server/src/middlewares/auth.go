package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cakezero/go-server/src/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

var IdKey = "user-id"
var ctx = context.Background()

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			utils.Response(res, "Missing Authorization header", "u")
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			utils.Response(res, "Invalid auth header format", "u")
			return
		}

		accessToken := authParts[1]

		revoked, redisErr := utils.GetRedisClient().Get(ctx, accessToken).Result()

		if redisErr != redis.Nil || revoked == "revoked" {
			utils.Response(res, "Access token has been revoked", "u")
			return
		}

		parsedToken, tokenParseErr := jwt.Parse(accessToken, func(token *jwt.Token) (interface {}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")	
			}

			return utils.JWT_SECRET, nil
		})

		if tokenParseErr != nil || !parsedToken.Valid {
			utils.Response(res, "Invalid or expired token", "u")
			return
		}

		if payload, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			id, ok := payload["id"].(string)
			if !ok {
				utils.Response(res, "token payload is invalid", "u")
				return
			}

			ctx := context.WithValue(req.Context(), IdKey, id)
			next.ServeHTTP(res, req.WithContext(ctx))
			return
		}

		utils.Response(res, "Unauthorized", "u")
	})
}

