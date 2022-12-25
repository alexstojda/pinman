package errors

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pinman/internal/app/generated"
)

func AbortWithError(code int, detail string, ctx *gin.Context, meta ...map[string]interface{}) {
	var metaVar *map[string]interface{}
	if len(meta) > 1 {
		metaVar = &meta[0]
	}

	ctx.AbortWithStatusJSON(code, &generated.ErrorResponse{
		Status: code,
		Title:  http.StatusText(code),
		Detail: detail,
		Meta:   metaVar,
	})
}
