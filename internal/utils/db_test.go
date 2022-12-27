package utils_test

import (
	g "github.com/onsi/ginkgo/v2"
	m "github.com/onsi/gomega"
	"pinman/internal/utils"
	"time"
)

var _ = g.Describe("AnyTime", func() {
	g.When("Match is called", func() {
		g.It("returns true if value is time.Time", func() {
			m.Expect(
				utils.AnyTime{}.Match(time.Now()),
			).To(m.BeTrue())
		})
		g.It("returns false if value is not time.Time", func() {
			m.Expect(
				utils.AnyTime{}.Match("time.Now()"),
			).To(m.BeFalse())
		})
	})
})

var _ = g.Describe("AnyString", func() {
	g.When("Match is called", func() {
		g.It("returns true if value is string", func() {
			m.Expect(
				utils.AnyString{}.Match("time.Now()"),
			).To(m.BeTrue())
		})
		g.It("returns false if value is not string", func() {
			m.Expect(
				utils.AnyString{}.Match(time.Now()),
			).To(m.BeFalse())
		})
	})
})
