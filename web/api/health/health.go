package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Health struct {
}

func NewHealth() *Health {
	return &Health{}
}

type Response struct {
	Status string `json:"status"`
}

func (h *Health) Get(c *gin.Context) {
	response := Response{Status: "ok"}
	c.JSON(http.StatusOK, response)
}
