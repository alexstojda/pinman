package utils_test

import (
	g "github.com/onsi/ginkgo/v2"
	m "github.com/onsi/gomega"
	"pinman/internal/utils"
)

var _ = g.Describe("NewGormMock", func() {
	g.When("is called", func() {
		g.It("should succeed", func() {
			db, mock := utils.NewGormMock()
			m.Expect(db).NotTo(m.BeNil())
			m.Expect(mock).NotTo(m.BeNil())
		})
	})
})

var _ = g.Describe("NewGinTestCtx", func() {
	g.When("is called", func() {
		g.It("should succeed", func() {
			ctx, responseRecorder, router := utils.NewGinTestCtx()
			m.Expect(ctx).NotTo(m.BeNil())
			m.Expect(responseRecorder).NotTo(m.BeNil())
			m.Expect(router).NotTo(m.BeNil())
		})
	})
})
