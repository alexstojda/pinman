package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"pinman/web/api/auth"
	"pinman/web/api/hello"
	"pinman/web/api/user"
	"pinman/web/health"
)

type Server struct {
	ClientOrigin string
	SPAPath      string
	GormDB       *gorm.DB
	Health       *health.Health
	Hello        *hello.Hello
	Auth         *auth.Controller
	User         *user.Controller
}

func NewServer(spaPath string, clientOrigin string, gormDb *gorm.DB) *Server {
	return &Server{
		ClientOrigin: clientOrigin,
		SPAPath:      spaPath,
		Health:       health.NewHealth(),
		Hello:        hello.NewHello(),
		Auth:         auth.NewController(gormDb),
		User:         user.NewController(gormDb),
	}
}

func (s *Server) StartServer() {
	router := gin.New()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", s.ClientOrigin}
	corsConfig.AllowCredentials = true

	router.Use(cors.New(corsConfig))

	router.Use(logger.SetLogger(logger.Config{
		SkipPath: []string{
			"/health",
			"/metrics",
		},
	}))

	router.Use(errorHandler)

	// API ROUTES
	apiRouter := router.Group("/api")

	apiRouter.GET("/health", s.Health.Get)
	apiRouter.GET("/hello", s.Hello.Get)

	auth.NewRouteController(s.Auth).AuthRoutes(apiRouter)
	user.NewRouteController(s.User).UserRoutes(apiRouter)

	// SPA ROUTE
	// Only loaded if SPAPath is defined.
	if s.SPAPath != "" {
		router.Use(static.Serve("/", static.LocalFile(s.SPAPath, true)))
	}

	// Uncomment below to enable prometheus metrics
	//ConfigurePrometheus(router, []string{})

	err := router.Run()
	if err != nil {
		log.Error().Msgf("Web server startup failed with error %s", err)
	}
}
