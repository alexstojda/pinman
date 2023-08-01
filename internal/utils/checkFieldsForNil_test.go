package utils_test

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"pinman/internal/utils"
)

type TestStruct struct {
	Field1 *string
	Field2 *string
	Field3 *string
}

var _ = ginkgo.Describe("CheckFieldsForNil", func() {
	ginkgo.When("given a struct with all fields set", func() {
		ginkgo.It("should return nil", func() {
			testStruct := TestStruct{
				Field1: utils.PtrString(""),
				Field2: utils.PtrString(""),
				Field3: utils.PtrString(""),
			}
			err := utils.CheckFieldsForNil(testStruct)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})
	ginkgo.When("given a struct with one field set", func() {
		ginkgo.It("should return an error", func() {
			testStruct := TestStruct{
				Field1: utils.PtrString(""),
				Field2: nil,
				Field3: nil,
			}
			err := utils.CheckFieldsForNil(testStruct)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
