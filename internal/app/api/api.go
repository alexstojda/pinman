package api

import (
	"github.com/gin-gonic/gin"
	"pinman/internal/app/api/league"
	"pinman/internal/app/api/user"
)

type AuthHandlers struct {
	Login   gin.HandlerFunc
	Refresh gin.HandlerFunc
}

// Server
// implements generated.ServerInterface
type Server struct {
	User   *user.Controller
	League *league.Controller
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

func (s *Server) PostLeagues(c *gin.Context) {
	s.League.CreateLeague(c)
}

func (s *Server) GetLeagues(c *gin.Context) {
	s.League.ListLeagues(c)
}

func (s *Server) GetLeaguesSlug(c *gin.Context, slug string) {
	s.League.GetLeagueWithSlug(c, slug)
}
