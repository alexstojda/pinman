package api

import (
	"github.com/gin-gonic/gin"
	"pinman/internal/app/api/auth"
	"pinman/internal/app/api/user"
)

type Server struct {
	Auth *auth.Controller
	User *user.Controller
}

func (s *Server) PostAuthLogin(c *gin.Context) {
	s.Auth.SignInUser(c)
}

func (s *Server) PostAuthRefresh(c *gin.Context) {
	s.Auth.RefreshAccessToken(c)
}

func (s *Server) PostAuthRegister(c *gin.Context) {
	s.Auth.SignUpUser(c)
}

func (s *Server) PostUsersMe(c *gin.Context) {
	s.User.GetMe(c)
}
