package middleware

import "github.com/gin-gonic/gin"

func NoCache() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "no-cache")
		return
	}
}
