package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"booking-event/internal/common/util"
	"booking-event/internal/modules/auth/model"
)

type AuthValidator interface {
	VerifyJWTToken(token string) (*model.TokenClaims, error)
}

func ExtractTokenFromBearer(token string) (string, error) {
	splitToken := strings.Split(token, "Bearer ")
	if len(splitToken) != 2 {
		return "", errors.New("invalid token")
	}
	return splitToken[1], nil
}

func AuthMiddleware(validator AuthValidator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.GetHeader("Authorization")
		if bearerToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		token, err := ExtractTokenFromBearer(bearerToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		claims, err := validator.VerifyJWTToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		ctx.Request = ctx.Request.WithContext(util.SetUserIDContext(ctx.Request.Context(), claims.UserID))
		ctx.Next()
	}
}

func AdminAuthMiddleware(validator AuthValidator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.GetHeader("Authorization")
		if bearerToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		token, err := ExtractTokenFromBearer(bearerToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		claims, err := validator.VerifyJWTToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		if !strings.EqualFold(string(claims.Role), string(model.RoleAdmin)) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		ctx.Request = ctx.Request.WithContext(util.SetUserIDContext(ctx.Request.Context(), claims.UserID))
		ctx.Next()
	}
}
