package hello_test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"pinman/internal/app/api/hello"
	"testing"

	g "github.com/onsi/ginkgo/v2"
	m "github.com/onsi/gomega"
)

func TestHello(t *testing.T) {
	m.RegisterFailHandler(g.Fail)
	g.RunSpecs(t, "Hello Suite")
}

var _ = g.Describe("Hello", func() {
	g.When("Get receives a request", func() {
		g.It("succeeds", func() {
			resp := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			ctx, router := gin.CreateTestContext(resp)

			controller := hello.NewHello()

			router.GET("/", controller.Get)

			ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			router.ServeHTTP(resp, ctx.Request)

			m.Expect(resp.Code).To(m.Equal(http.StatusOK))

			data := hello.Response{}

			err := json.Unmarshal(resp.Body.Bytes(), &data)
			m.Expect(err).To(m.BeNil())
			m.Expect(data).To(m.Equal(hello.Response{Hello: "world"}))
		})
	})
})
