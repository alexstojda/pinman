package app

import "github.com/gin-gonic/gin"

func errorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		c.JSON(500, c.Errors)
	}
}
