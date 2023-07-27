package utils_test

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"pinman/internal/utils"
)

var _ = ginkgo.Describe("Slugify", func() {
	ginkgo.Context("when it receives a string to be slugified", func() {
		ginkgo.It("should return a slugified string", func() {
			value := "This is a string to be slugified"
			expected := "this-is-a-string-to-be-slugified"
			gomega.Expect(utils.Slugify(value)).To(gomega.Equal(expected))
		})
	})
	ginkgo.Context("when it receives a string with special characters to be slugified", func() {
		ginkgo.It("should return a slugified string", func() {
			value := "This is a string to be slugified with special characters: áéíóúñ"
			expected := "this-is-a-string-to-be-slugified-with-special-characters-aeioun"
			gomega.Expect(utils.Slugify(value)).To(gomega.Equal(expected))
		})
	})
	ginkgo.Context("when it receives a string to be slugified and a max length", func() {
		ginkgo.It("should return a slugified string with the max length", func() {
			value := "This is a string to be slugified"
			expected := "this-is-a-string"
			gomega.Expect(utils.Slugify(value, 17)).To(gomega.Equal(expected))
		})
	})
})
