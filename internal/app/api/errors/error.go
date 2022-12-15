package errors

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Error struct {
	Status string `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Meta   any    `json:"meta,omitempty"`
}

func AbortWithError(code int, detail string, ctx *gin.Context, meta ...any) {
	var metaVar any
	if len(meta) > 1 {
		metaVar = meta[0]
	}

	ctx.AbortWithStatusJSON(code, &Error{
		Status: strconv.FormatInt(int64(code), 10),
		Title:  http.StatusText(code),
		Detail: detail,
		Meta:   metaVar,
	})
}
