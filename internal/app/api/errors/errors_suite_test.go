package errors_test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"pinman/internal/app/api/errors"
	"pinman/internal/app/generated"
	"pinman/internal/utils"
	"testing"
	"time"

	g "github.com/onsi/ginkgo/v2"
	m "github.com/onsi/gomega"
)

func TestErrors(t *testing.T) {
	m.RegisterFailHandler(g.Fail)
	g.RunSpecs(t, "Errors Suite")
}

var _ = g.Describe("AbortWithError", func() {
	var ctx *gin.Context
	var rr *httptest.ResponseRecorder

	g.BeforeEach(func() {
		ctx, rr, _ = utils.NewGinTestCtx()
	})
	g.When("called with no metadata", func() {
		g.It("should return correct response", func() {
			errors.AbortWithError(http.StatusNotFound, "page does not exist", ctx)

			m.Expect(rr.Code).To(m.BeEquivalentTo(http.StatusNotFound))

			expected := &generated.ErrorResponse{
				Detail: "page does not exist",
				Status: http.StatusNotFound,
				Title:  http.StatusText(http.StatusNotFound),
			}

			response := &generated.ErrorResponse{}
			err := json.Unmarshal(rr.Body.Bytes(), response)
			m.Expect(err).To(m.BeNil())
			m.Expect(response).To(m.Equal(expected))
		})
	})
	g.When("called with metadata", func() {
		g.It("should return correct response", func() {
			meta := map[string]interface{}{
				"foo":  "bar",
				"time": time.Now().Format(time.RFC3339),
			}

			errors.AbortWithError(http.StatusNotFound, "page does not exist", ctx, meta)

			m.Expect(rr.Code).To(m.BeEquivalentTo(http.StatusNotFound))

			expected := &generated.ErrorResponse{
				Detail: "page does not exist",
				Status: http.StatusNotFound,
				Title:  http.StatusText(http.StatusNotFound),
				Meta:   &meta,
			}

			response := &generated.ErrorResponse{}
			err := json.Unmarshal(rr.Body.Bytes(), response)
			m.Expect(err).To(m.BeNil())
			m.Expect(response).To(m.Equal(expected))
		})
	})
})
