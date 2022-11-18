package auth

import (
	"github.com/gin-gonic/gin"
	"pinman/web/middleware"
)

type RouteController struct {
	authController *Controller
}

func NewRouteController(authController *Controller) *RouteController {
	return &RouteController{authController}
}

func (rc *RouteController) AuthRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/auth")

	router.POST("/register", rc.authController.SignUpUser)
	router.POST("/login", rc.authController.SignInUser)
	router.GET("/refresh", rc.authController.RefreshAccessToken)
	router.GET("/logout", middleware.DeserializeUser(), rc.authController.LogoutUser)
}
