package health_test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"pinman/internal/app/api/health"
	"testing"

	g "github.com/onsi/ginkgo/v2"
	m "github.com/onsi/gomega"
)

func TestHealth(t *testing.T) {
	m.RegisterFailHandler(g.Fail)
	g.RunSpecs(t, "Health Suite")
}

var _ = g.Describe("Health", func() {
	g.When("Get receives a request", func() {
		g.It("succeeds", func() {
			resp := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			ctx, router := gin.CreateTestContext(resp)

			controller := health.NewHealth()

			router.GET("/", controller.Get)

			ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			router.ServeHTTP(resp, ctx.Request)

			m.Expect(resp.Code).To(m.Equal(http.StatusOK))

			data := health.Response{}

			err := json.Unmarshal(resp.Body.Bytes(), &data)
			m.Expect(err).To(m.BeNil())
			m.Expect(data).To(m.Equal(health.Response{Status: "ok"}))
		})
	})
})
