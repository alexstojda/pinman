package api

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pinman/internal/app/api/league"
	"pinman/internal/app/api/location"
	"pinman/internal/app/api/tournament"
	"pinman/internal/app/api/user"
	"pinman/internal/utils"
)

type AuthHandlers struct {
	Login   gin.HandlerFunc
	Refresh gin.HandlerFunc
}

// Server
// implements generated.ServerInterface
type Server struct {
	User       *user.Controller
	League     *league.Controller
	Location   *location.Controller
	Tournament *tournament.Controller
	AuthHandlers
}

func NewServer(db *gorm.DB, authMiddleware *jwt.GinJWTMiddleware) *Server {
	server := &Server{
		User:       user.NewController(db),
		League:     league.NewController(db),
		Location:   location.NewController(db),
		Tournament: tournament.NewController(db),
		AuthHandlers: AuthHandlers{
			Login:   authMiddleware.LoginHandler,
			Refresh: authMiddleware.RefreshHandler,
		},
	}
	if err := utils.CheckFieldsForNil(server); err != nil {
		panic(fmt.Errorf("failed to create server: %w", err))
	}

	return server
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

func (s *Server) GetLocations(c *gin.Context) {
	s.Location.ListLocations(c)
}

func (s *Server) PostLocations(c *gin.Context) {
	s.Location.CreateLocation(c)
}

func (s *Server) GetLocationsSlug(c *gin.Context, slug string) {
	s.Location.GetLocationWithSlug(c, slug)
}

func (s *Server) PostTournaments(c *gin.Context) {
	s.Tournament.CreateTournament(c)
}

func (s *Server) GetTournaments(c *gin.Context) {
	s.Tournament.ListTournaments(c)
}
