package pinballmap_test

import (
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/slices"
	"net/http"
	"pinman/internal/clients/generic"
	"pinman/internal/clients/pinballmap"
	"reflect"
	"testing"
)

func TestPinballMap(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "PinballMap Suite")
}

var _ = ginkgo.Describe("Client", func() {
	var mgClient *generic.MockClientInterface

	ginkgo.BeforeEach(func() {
		mgClient = generic.NewMockClientInterface(ginkgo.GinkgoT())
	})

	ginkgo.When("a client is created", func() {
		ginkgo.It("returns a new Client", func() {
			client := pinballmap.NewClient()
			gomega.Expect(client).ToNot(gomega.BeNil())
		})
	})

	// Test GetLocations
	ginkgo.When("GetLocations is called", func() {
		ginkgo.Context("using the live API with 'north star' filter", func() {
			ginkgo.It("should return a list of locations that contains at least an entry for 'North Star Machines à Piastres'", func() {
				client := pinballmap.NewClient()
				locations, err := client.GetLocations("north star")
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(locations).ToNot(gomega.BeNil())
				gomega.Expect(len(locations)).To(gomega.BeNumerically(">", 0))

				slices.ContainsFunc(locations, func(location pinballmap.Location) bool {
					return location.Name == "North Star Machines à Piastres"
				})
			})
		})
		ginkgo.Context("and the API call fails with unknown error", func() {
			ginkgo.It("should return an error", func() {
				mgClient.On(
					"Do",
					mock.AnythingOfType(reflect.TypeOf(&http.Request{}).String()),
					&pinballmap.LocationsResponse{},
					&pinballmap.ErrorResponse{},
				).Return(-1, fmt.Errorf("some error"))

				client := pinballmap.NewClientWithGenericClient(mgClient)
				locations, err := client.GetLocations("north star")

				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(locations).To(gomega.BeNil())
			})
		})
		ginkgo.Context("and the API returns an error response", func() {
			ginkgo.It("should return an error", func() {
				mgClient.On(
					"Do",
					mock.AnythingOfType(reflect.TypeOf(&http.Request{}).String()),
					&pinballmap.LocationsResponse{},
					&pinballmap.ErrorResponse{},
				).Run(func(args mock.Arguments) {
					errResponse := args.Get(2).(*pinballmap.ErrorResponse)
					errResponse.Errors = "some error"
				}).Return(500, nil)

				client := pinballmap.NewClientWithGenericClient(mgClient)
				locations, err := client.GetLocations("north star")

				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("some error"))
				gomega.Expect(locations).To(gomega.BeNil())
			})
		})
	})

	// Test GetLocation
	ginkgo.When("GetLocation is called", func() {
		ginkgo.Context("using the live API with location ID 7464", func() {
			ginkgo.It("should return a location with ID 7464", func() {
				client := pinballmap.NewClient()
				location, err := client.GetLocation(7464)
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(location).ToNot(gomega.BeNil())
				gomega.Expect(location.ID).To(gomega.Equal(7464))
			})
		})
		ginkgo.Context("and the API returns an error", func() {
			ginkgo.It("should return an error", func() {
				mgClient.On(
					"Do",
					mock.AnythingOfType(reflect.TypeOf(&http.Request{}).String()),
					&pinballmap.Location{},
					&pinballmap.ErrorResponse{},
				).Run(func(args mock.Arguments) {
					errResponse := args.Get(2).(*pinballmap.ErrorResponse)
					errResponse.Errors = "some error"
				}).Return(-1, nil)

				client := pinballmap.NewClientWithGenericClient(mgClient)
				location, err := client.GetLocation(1)

				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("some error"))
				gomega.Expect(location).To(gomega.BeNil())
			})
		})
		ginkgo.Context("and performing the request returns an error", func() {
			ginkgo.It("should return an error", func() {
				mgClient.On(
					"Do",
					mock.AnythingOfType(reflect.TypeOf(&http.Request{}).String()),
					&pinballmap.Location{},
					&pinballmap.ErrorResponse{},
				).Return(-1, fmt.Errorf("some error"))

				client := pinballmap.NewClientWithGenericClient(mgClient)
				location, err := client.GetLocation(1)

				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(location).To(gomega.BeNil())
			})
		})
	})
})
