package api

import (
	"github.com/gin-gonic/gin"
	"pinman/internal/app/api/user"
)

type AuthHandlers struct {
	Login   gin.HandlerFunc
	Refresh gin.HandlerFunc
}

// Server
// implements generated.ServerInterface
type Server struct {
	User *user.Controller
	AuthHandlers
}

func (s *Server) PostAuthLogin(c *gin.Context) {
	s.AuthHandlers.Login(c)
}

func (s *Server) GetAuthRefresh(c *gin.Context) {
	s.AuthHandlers.Refresh(c)
}

func (s *Server) PostUsersRegister(c *gin.Context) {
	s.User.SignUpUser(c)
}

func (s *Server) GetUsersMe(c *gin.Context) {
	s.User.GetMe(c)
}
