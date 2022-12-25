package app

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"pinman/internal/app/api"
	"pinman/internal/app/api/auth"
	"pinman/internal/app/api/errors"
	"pinman/internal/app/api/hello"
	"pinman/internal/app/api/user"
	"pinman/internal/app/generated"
	"pinman/internal/app/health"
	"pinman/internal/app/middleware"
	"pinman/internal/utils"
)

type Server struct {
	SPACacheDisabled bool
	Db               *gorm.DB
	Health           *health.Health
	Hello            *hello.Hello
	Config           *utils.Config
}

func NewServer(config *utils.Config, db *gorm.DB) *Server {
	return &Server{
		Health: health.NewHealth(),
		Hello:  hello.NewHello(),
		Db:     db,
		Config: config,
	}
}

func (s *Server) StartServer() error {
	router := gin.New()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = s.Config.ClientOrigins
	if len(corsConfig.AllowOrigins) > 0 {
		log.Info().Interface("allowedOrigins", corsConfig.AllowOrigins).Msg("CORS origins configured")
	}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = append(
		corsConfig.AllowHeaders,
		[]string{
			"Authorization",
		}...,
	)

	router.Use(cors.New(corsConfig))

	// Since we don't use any proxy, this feature can be disabled
	err := router.SetTrustedProxies(nil)
	if err != nil {
		return fmt.Errorf("could not set trusted proxies: %v", err)
	}

	router.Use(logger.SetLogger(logger.Config{
		SkipPath: []string{
			"/health",
			"/metrics",
		},
	}))

	authMiddleware, err := auth.CreateJWTMiddleware(s.Config, s.Db)
	if err != nil {
		return fmt.Errorf("could not initialize JWT middleware: %v", err)
	}

	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		errors.AbortWithError(404, "page not found", c)
		return
	})

	generated.RegisterHandlersWithOptions(
		router,
		&api.Server{
			User: user.NewController(s.Db),
			AuthHandlers: api.AuthHandlers{
				Login:   authMiddleware.LoginHandler,
				Refresh: authMiddleware.RefreshHandler,
			},
		},
		generated.GinServerOptions{
			BaseURL: "/api",
			Middlewares: []generated.MiddlewareFunc{
				func(c *gin.Context) {
					if _, ok := c.Get(generated.PinmanAuthScopes); ok {
						authMiddleware.MiddlewareFunc()(c)
					}
				},
			},
			ErrorHandler: nil,
		})

	router.GET("/api/health", s.Health.Get)
	router.GET("/api/hello", s.Hello.Get)

	// SPA ROUTE
	// Only loaded if SPAPath is defined.
	if s.Config.SPAPath != "" {
		log.Debug().Str("spaPath", s.Config.SPAPath).Msg("SPA_PATH is set, will serve")

		spaRoute := static.Serve("/", static.LocalFile(s.Config.SPAPath, true))

		if s.SPACacheDisabled {
			router.Use(middleware.NoCache()).Use(spaRoute)
		} else {
			router.Use(spaRoute)
		}
	}

	// Uncomment below to enable prometheus metrics
	//ConfigurePrometheus(router, []string{})

	return router.Run()
}
