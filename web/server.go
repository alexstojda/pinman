package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"gorm.io/gorm"
	"pinman/web/api/auth"
	"pinman/web/api/health"
	"pinman/web/api/hello"
	"pinman/web/api/user"
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

	prometheus := ginprometheus.NewPrometheus("gin")

	// Prevents high cardinality of metrics Source: https://github.com/zsais/go-gin-prometheus#preserving-a-low-cardinality-for-the-request-counter
	prometheus.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.Path // Query params are dropped here so there is not a metric for every permutation of query param usage on a route

		//  If a route uses parameters, replace the parameter value with its name. Else there will be a metric for the route
		//  with every possible value of that parameter and this will cause performance issues in Prometheus.
		//
		//  If your service uses route parameters, uncomment the for loop below and add a case for each parameter. The example case
		//  below works for routes with a parameter called 'name', like '/api/function/:name'
		//  --
		//    for _, p := range c.Params {
		//      switch p.Key {
		//      case "name":
		//        url = strings.Replace(url, p.Value, ":name", 1)
		//      }
		//    }
		return url
	}
	prometheus.Use(router)

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

	err := router.Run()
	if err != nil {
		log.Error().Msgf("Web server startup failed with error %s", err)
	}
}
