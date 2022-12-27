package utils_test

import (
	g "github.com/onsi/ginkgo/v2"
	m "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
	"pinman/internal/utils"
)

var _ = g.Describe("password.go", func() {
	g.When("HashPassword is called", func() {
		g.Context("with invalid HashCost", func() {
			g.It("should return an error", func() {
				utils.HashCost = 100
				str, err := utils.HashPassword("password")
				m.Expect(err).NotTo(m.BeNil())
				m.Expect(str).To(m.BeEmpty())
			})
		})

		g.Context("with valid HashCost", func() {
			g.It("should return the string, hashed", func() {
				val := "password"
				str, err := utils.HashPassword(val)
				m.Expect(err).To(m.BeNil())
				m.Expect(utils.VerifyPassword(str, val)).To(m.BeNil())
			})
		})
	})
	g.AfterEach(func() {
		utils.HashCost = bcrypt.DefaultCost
	})
})
