package pinballmap_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/slices"
	"io"
	"net/http"
	"net/url"
	"pinman/internal/clients/pinballmap"
	"pinman/internal/utils"
	"reflect"
	"testing"
)

func TestPinballMap(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "PinballMap Suite")
}

var _ = ginkgo.Describe("Client", func() {
	var mockDoer *utils.MockHttpDoer

	ginkgo.BeforeEach(func() {
		mockDoer = utils.NewMockHttpDoer(ginkgo.GinkgoT())
	})

	ginkgo.When("a client is created", func() {
		var urrl *url.URL
		ginkgo.BeforeEach(func() {
			client := pinballmap.NewClient()
			urrl = client.NewUrl("/api/v1/locations.json", nil)
		})
		ginkgo.It("should have a default api scheme of https", func() {
			gomega.Expect(urrl.Scheme).To(gomega.Equal("https"))
		})
		ginkgo.It("should have a default api host of pinballmap.com", func() {
			gomega.Expect(urrl.Host).To(gomega.Equal("pinballmap.com"))
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
		ginkgo.Context("and the API returns an error of unknown structure", func() {
			ginkgo.It("should return an error", func() {
				mockDoer.On(
					"Do",
					mock.AnythingOfType(reflect.TypeOf(&http.Request{}).String()),
				).Return(&http.Response{
					StatusCode: 500,
					Status:     "Internal Server Error",
					Body:       io.NopCloser(bytes.NewReader([]byte("some error"))),
				}, nil)

				client := pinballmap.NewClientWithHttpClient(mockDoer)
				locations, err := client.GetLocations("north star")

				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(locations).To(gomega.BeNil())
			})
		})
		ginkgo.Context("and the API returns an error", func() {
			ginkgo.It("should return an error", func() {
				mockError, err := json.Marshal(pinballmap.ErrorResponse{
					Errors: "some error",
				})
				gomega.Expect(err).To(gomega.BeNil())

				mockDoer.On(
					"Do",
					mock.AnythingOfType(reflect.TypeOf(&http.Request{}).String()),
				).Return(&http.Response{
					StatusCode: 500,
					Status:     "Internal Server Error",
					Body:       io.NopCloser(bytes.NewReader(mockError)),
				}, nil)

				client := pinballmap.NewClientWithHttpClient(mockDoer)
				locations, err := client.GetLocations("north star")

				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("some error"))
				gomega.Expect(locations).To(gomega.BeNil())
			})
		})
		ginkgo.Context("and performing the request returns an error", func() {
			ginkgo.It("should return an error", func() {
				mockDoer.On(
					"Do",
					mock.AnythingOfType(reflect.TypeOf(&http.Request{}).String()),
				).Return(nil, fmt.Errorf("some error"))

				client := pinballmap.NewClientWithHttpClient(mockDoer)
				locations, err := client.GetLocations("north star")

				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(locations).To(gomega.BeNil())
			})
		})
		ginkgo.Context("and response is not the expected structure", func() {
			ginkgo.It("should return an error", func() {
				mockDoer.On(
					"Do",
					mock.AnythingOfType(reflect.TypeOf(&http.Request{}).String()),
				).Return(&http.Response{
					Status:     http.StatusText(http.StatusOK),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
				}, nil)

				client := pinballmap.NewClientWithHttpClient(mockDoer)
				locations, err := client.GetLocations("north star")

				gomega.Expect(err).To(gomega.HaveOccurred())
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
				mockError, err := json.Marshal(pinballmap.ErrorResponse{
					Errors: "some error",
				})
				gomega.Expect(err).To(gomega.BeNil())

				mockDoer.On(
					"Do",
					mock.AnythingOfType(reflect.TypeOf(&http.Request{}).String()),
				).Return(&http.Response{
					StatusCode: 500,
					Status:     "Internal Server Error",
					Body:       io.NopCloser(bytes.NewReader(mockError)),
				}, nil)

				client := pinballmap.NewClientWithHttpClient(mockDoer)
				location, err := client.GetLocation(1)

				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("some error"))
				gomega.Expect(location).To(gomega.BeNil())
			})
		})
		ginkgo.Context("and performing the request returns an error", func() {
			ginkgo.It("should return an error", func() {
				mockDoer.On(
					"Do",
					mock.AnythingOfType(reflect.TypeOf(&http.Request{}).String()),
				).Return(nil, fmt.Errorf("some error"))

				client := pinballmap.NewClientWithHttpClient(mockDoer)
				location, err := client.GetLocation(1)

				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(location).To(gomega.BeNil())
			})
		})
		ginkgo.Context("and the API returns an error of unknown structure", func() {
			ginkgo.It("should return an error", func() {
				mockDoer.On(
					"Do",
					mock.AnythingOfType(reflect.TypeOf(&http.Request{}).String()),
				).Return(&http.Response{
					StatusCode: 500,
					Status:     "Internal Server Error",
					Body:       io.NopCloser(bytes.NewReader([]byte("some error"))),
				}, nil)

				client := pinballmap.NewClientWithHttpClient(mockDoer)
				location, err := client.GetLocation(1)

				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(location).To(gomega.BeNil())
			})
		})
		ginkgo.Context("and response is not the expected structure", func() {
			ginkgo.It("should return an error", func() {
				mockDoer.On(
					"Do",
					mock.AnythingOfType(reflect.TypeOf(&http.Request{}).String()),
				).Return(&http.Response{
					Status:     http.StatusText(http.StatusOK),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
				}, nil)

				client := pinballmap.NewClientWithHttpClient(mockDoer)
				location, err := client.GetLocation(1)

				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(location).To(gomega.BeNil())
			})
		})
	})
})
