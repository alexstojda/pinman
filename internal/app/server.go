package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	auth2 "pinman/internal/app/api/auth"
	"pinman/internal/app/api/hello"
	user2 "pinman/internal/app/api/user"
	"pinman/internal/app/health"
)

type Server struct {
	ClientOrigins []string
	SPAPath       string
	GormDB        *gorm.DB
	Health        *health.Health
	Hello         *hello.Hello
	Auth          *auth2.Controller
	User          *user2.Controller
}

func NewServer(spaPath string, clientOrigins []string, gormDb *gorm.DB) *Server {
	return &Server{
		ClientOrigins: clientOrigins,
		SPAPath:       spaPath,
		Health:        health.NewHealth(),
		Hello:         hello.NewHello(),
		Auth:          auth2.NewController(gormDb),
		User:          user2.NewController(gormDb),
	}
}

func (s *Server) StartServer() {
	router := gin.New()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = s.ClientOrigins
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	// Since we don't use any proxy, this feature can be disabled
	err := router.SetTrustedProxies(nil)
	if err != nil {
		return
	}

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

	auth2.NewRouteController(s.Auth).AuthRoutes(apiRouter)
	user2.NewRouteController(s.User).UserRoutes(apiRouter)

	// SPA ROUTE
	// Only loaded if SPAPath is defined.
	log.Debug().Str("spaPath", s.SPAPath).Msg("SPA_PATH is set, will serve")
	if s.SPAPath != "" {
		router.Use(static.Serve("/", static.LocalFile(s.SPAPath, true)))
	}

	// Uncomment below to enable prometheus metrics
	//ConfigurePrometheus(router, []string{})

	err = router.Run()
	if err != nil {
		log.Error().Msgf("Web server startup failed with error %s", err)
	}
}
