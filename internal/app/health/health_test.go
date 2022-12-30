package health

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth_Get(t *testing.T) {

	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, router := gin.CreateTestContext(resp)

	health := NewHealth()

	router.GET("/test", health.Get)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(resp, ctx.Request)

	assert.Equal(t, resp.Code, 200)

	data := Response{}

	err := json.Unmarshal(resp.Body.Bytes(), &data)
	if err != nil {
		t.Log(err.Error())
		t.Fatalf("Could not unmarshal response body. Got %s", resp.Body.String())
	}

	assert.Equal(t, data, Response{Status: "ok"})
}
