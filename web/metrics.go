package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"strings"
)

// ConfigurePrometheus configures gin to collect metrics for prometheus to scrape
// router *gin.Engine gin instance to configure and enable prometheus for
// paramNames []string a complete list of route parameter names used by the application (see note below)
// ---
// Route Parameters
//
// If a route uses parameters, we must replace the parameter value with its name. Otherwise, there will be a metric for
// the route with every possible value of that parameter and this will cause performance issues in Prometheus.
//
// Example: If the application has a route with a parameter called 'name', like '/api/function/:name',
// add `"name"` to `paramNames`
func ConfigurePrometheus(router *gin.Engine, paramNames []string) {
	prometheus := ginprometheus.NewPrometheus("gin")

	// Prevents high cardinality of metrics Source: https://github.com/zsais/go-gin-prometheus#preserving-a-low-cardinality-for-the-request-counter
	prometheus.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.Path // Query params are dropped here so there is not a metric for every permutation of query param usage on a route

		for _, urlParam := range c.Params {
			for _, knownParam := range paramNames {
				if urlParam.Key == knownParam {
					url = strings.Replace(url, urlParam.Value, fmt.Sprintf(":%s", knownParam), 1)
					break
				}
			}
		}
		return url
	}
	prometheus.Use(router)
}
