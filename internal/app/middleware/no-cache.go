package middleware

import (
	"github.com/gin-gonic/gin"
)

func NoCache() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "no-cache")
		return
	}
}
