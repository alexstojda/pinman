package user

import (
	"github.com/gin-gonic/gin"
	"pinman/web/middleware"
)

type RouteController struct {
	userController *Controller
}

func NewRouteController(userController *Controller) *RouteController {
	return &RouteController{userController}
}

func (uc *RouteController) UserRoutes(rg *gin.RouterGroup) {

	router := rg.Group("users")
	router.GET("/me", middleware.DeserializeUser(), uc.userController.GetMe)
}
