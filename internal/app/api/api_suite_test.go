package api_test

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"gorm.io/gorm"
	"pinman/internal/app/api"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestApi(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Api Suite")
}

var _ = ginkgo.Describe("NewServer", func() {
	ginkgo.It("should not return nil", func() {
		server := api.NewServer(
			&gorm.DB{},
			&jwt.GinJWTMiddleware{},
		)
		gomega.Expect(server).NotTo(gomega.BeNil())
	})
})
