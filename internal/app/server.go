package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"pinman/internal/app/api"
	"pinman/internal/app/api/auth"
	"pinman/internal/app/api/hello"
	"pinman/internal/app/api/user"
	"pinman/internal/app/generated"
	"pinman/internal/app/health"
	"pinman/internal/app/middleware"
)

type Server struct {
	ClientOrigins    []string
	SPAPath          string
	SPACacheDisabled bool
	Db               *gorm.DB
	Health           *health.Health
	Hello            *hello.Hello
	ApiServer        generated.ServerInterface
}

func NewServer(spaPath string, clientOrigins []string, db *gorm.DB) *Server {
	return &Server{
		ClientOrigins: clientOrigins,
		SPAPath:       spaPath,
		Health:        health.NewHealth(),
		Hello:         hello.NewHello(),
		ApiServer: &api.Server{
			Auth: auth.NewController(db),
			User: user.NewController(db),
		},
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

	generated.RegisterHandlersWithOptions(router, s.ApiServer, generated.GinServerOptions{
		BaseURL: "/api",
		Middlewares: []generated.MiddlewareFunc{
			middleware.AuthenticateUser(s.Db),
		},
		ErrorHandler: nil,
	})

	router.GET("/api/health", s.Health.Get)
	router.GET("/api/hello", s.Hello.Get)

	// SPA ROUTE
	// Only loaded if SPAPath is defined.
	if s.SPAPath != "" {
		log.Debug().Str("spaPath", s.SPAPath).Msg("SPA_PATH is set, will serve")

		spaRoute := static.Serve("/", static.LocalFile(s.SPAPath, true))

		if s.SPACacheDisabled {
			router.Use(middleware.NoCache()).Use(spaRoute)
		} else {
			router.Use(spaRoute)
		}
	}

	// Uncomment below to enable prometheus metrics
	//ConfigurePrometheus(router, []string{})

	err = router.Run()
	if err != nil {
		log.Error().Msgf("Web server startup failed with error %s", err)
	}
}
