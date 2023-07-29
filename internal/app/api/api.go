package api

import (
	"errors"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pinman/internal/app/api/league"
	"pinman/internal/app/api/location"
	"pinman/internal/app/api/user"
	"reflect"
)

type AuthHandlers struct {
	Login   gin.HandlerFunc
	Refresh gin.HandlerFunc
}

// Server
// implements generated.ServerInterface
type Server struct {
	User     *user.Controller
	League   *league.Controller
	Location *location.Controller
	AuthHandlers
}

func NewServer(db *gorm.DB, authMiddleware *jwt.GinJWTMiddleware) *Server {
	server := &Server{
		User:     user.NewController(db),
		League:   league.NewController(db),
		Location: location.NewController(db),
		AuthHandlers: AuthHandlers{
			Login:   authMiddleware.LoginHandler,
			Refresh: authMiddleware.RefreshHandler,
		},
	}
	if err := checkFieldsForNil(server); err != nil {
		panic(fmt.Errorf("failed to create server: %w", err))
	}

	return server
}

func checkFieldsForNil(s *Server) error {
	serverValue := reflect.ValueOf(*s)
	for i := 0; i < serverValue.NumField(); i++ {
		fieldValue := serverValue.Field(i)
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			return errors.New(fmt.Sprint("field ", serverValue.Type().Field(i).Name, " is nil"))
		}
	}
	return nil
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
