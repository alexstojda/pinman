package web

import (
	"pinman/web/api/health"
	"pinman/web/api/hello"

	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

type Server struct {
	SPAPath string
	Health  *health.Health
	Hello   *hello.Hello
}

func NewServer(spaPath string) *Server {
	return &Server{
		SPAPath: spaPath,
		Health:  health.NewHealth(),
		Hello:   hello.NewHello(),
	}
}

func (s *Server) StartServer() {
	router := gin.New()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

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

		//  If a route uses parameters, replace the parameter value with its name. Else there will be a metric for the route with
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
	router.GET("/api/health", s.Health.Get)
	router.GET("/api/hello", s.Hello.Get)

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
