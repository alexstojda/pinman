package utils_test

import (
	"pinman/internal/utils"
	"testing"
	"time"

	g "github.com/onsi/ginkgo/v2"
	m "github.com/onsi/gomega"
)

func TestUtils(t *testing.T) {
	m.RegisterFailHandler(g.Fail)
	g.RunSpecs(t, "Utils Suite")
}

var _ = g.Describe("utils.go", func() {
	g.When("PtrString is called", func() {
		g.It("returns a pointer to a string with same value", func() {
			val := "value"
			ptrVal := utils.PtrString(val)

			m.Expect(ptrVal).To(m.BeAssignableToTypeOf(&val))
			m.Expect(ptrVal).To(m.BeEquivalentTo(&val))
		})
	})

	g.When("FormatTime is called", func() {
		g.It("returns the time formatted in RFC3339", func() {
			val, _ := time.Parse(time.RFC3339Nano, time.Now().Format(time.RFC3339Nano))
			strResult := utils.FormatTime(val)
			timeResult, err := time.Parse(time.RFC3339Nano, strResult)

			m.Expect(err).To(m.BeNil())
			m.Expect(timeResult).To(m.BeEquivalentTo(val))
		})
	})
})
