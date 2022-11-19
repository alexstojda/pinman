package hello

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Hello struct {
}

func NewHello() *Hello {
	return &Hello{}
}

type Response struct {
	Hello string `json:"hello"`
}

func (h *Hello) Get(c *gin.Context) {
	response := Response{
		Hello: "world",
	}
	c.JSON(http.StatusOK, response)
}
