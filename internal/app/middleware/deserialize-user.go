package middleware

import (
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"pinman/internal/app/generated"
	"pinman/internal/models"
	"pinman/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthenticateUser(db *gorm.DB) generated.MiddlewareFunc {
	return func(ctx *gin.Context) {
		// If scopes are not set, this route does not require authentication.
		if _, exists := ctx.Get(generated.OauthScopes); !exists {
			return
		}

		var accessToken string
		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err == nil {
			accessToken = cookie
		}

		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		config, _ := utils.LoadConfig(".")
		sub, err := utils.ValidateToken(accessToken, config.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		var user models.User
		result := db.First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
