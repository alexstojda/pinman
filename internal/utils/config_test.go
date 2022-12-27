package utils_test

import (
	g "github.com/onsi/ginkgo/v2"
	m "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"os"
	"pinman/internal/utils"
)

var _ = g.Describe("LoadConfig", func() {
	g.When("is called and ENV_FILE is set", func() {
		g.It("should succeed", func() {
			_ = os.Setenv("ENV_FILE", "../../.env")
			cfg, err := utils.LoadConfig()
			m.Expect(err).Error().To(m.BeNil())
			m.Expect(cfg).ToNot(m.BeNil())
		})
	})
	g.When("is called and ENV_FILE is not set", func() {
		g.It("should use default ENV_FILE and fail", func() {
			cfg, err := utils.LoadConfig()
			m.Expect(err).ToNot(m.BeNil())
			m.Expect(cfg).To(m.BeNil())
		})
	})
	g.AfterEach(func() {
		viper.Reset()
		_ = os.Unsetenv("ENV_FILE")
	})
})
