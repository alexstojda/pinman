package generic_test

import (
	"bytes"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/url"
	"pinman/internal/clients/generic"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestGeneric(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Generic Suite")
}

var _ = ginkgo.Describe("FormatRequestUrl", func() {
	ginkgo.When("is called with a host and path", func() {
		ginkgo.It("returns a formatted url", func() {
			url := generic.FormatRequestUrl("https://example.com", "/foo")
			gomega.Expect(url).To(gomega.Equal("https://example.com/foo"))
		})
	})
	ginkgo.When("is called with a host, path, and params", func() {
		ginkgo.It("returns a formatted url", func() {
			params := &url.Values{}
			params.Add("foo", "bar")
			url := generic.FormatRequestUrl("https://example.com", "/foo", params)
			gomega.Expect(url).To(gomega.Equal("https://example.com/foo?foo=bar"))
		})
	})
	ginkgo.When("is called with a host, path, and params with multiple values", func() {
		ginkgo.It("returns a formatted url", func() {
			params := &url.Values{}
			params.Add("foo", "bar")
			params.Add("foo", "baz")
			url := generic.FormatRequestUrl("https://example.com", "/foo", params)
			gomega.Expect(url).To(gomega.Equal("https://example.com/foo?foo=bar&foo=baz"))
		})
	})
	ginkgo.When("is called with a host that ends with a slash", func() {
		ginkgo.It("returns a formatted url", func() {
			url := generic.FormatRequestUrl("https://example.com/", "/foo")
			gomega.Expect(url).To(gomega.Equal("https://example.com/foo"))
		})
	})
	ginkgo.When("is called with a path that does not start with a slash", func() {
		ginkgo.It("returns a formatted url", func() {
			url := generic.FormatRequestUrl("https://example.com", "foo")
			gomega.Expect(url).To(gomega.Equal("https://example.com/foo"))
		})
	})
})

var _ = ginkgo.Describe("Client", func() {
	ginkgo.When("a client is created", func() {
		ginkgo.It("returns a new Client", func() {
			client := generic.NewClient()
			gomega.Expect(client).ToNot(gomega.BeNil())
		})
	})

	ginkgo.When("a client is created with a custom http client", func() {
		ginkgo.It("returns a new Client", func() {
			client := generic.NewClientWithHttpClient(nil)
			gomega.Expect(client).ToNot(gomega.BeNil())
		})
	})

	ginkgo.Describe("Do", func() {
		var mockHttpClient *generic.MockHttpDoer
		type Response struct {
			Foo string `json:"foo"`
		}

		ginkgo.BeforeEach(func() {
			mockHttpClient = generic.NewMockHttpDoer(ginkgo.GinkgoT())
		})

		ginkgo.When("is called with a nil response object", func() {
			ginkgo.It("returns an error", func() {
				client := generic.NewClientWithHttpClient(mockHttpClient)
				_, err := client.Do(nil, nil, nil)
				gomega.Expect(err).ToNot(gomega.BeNil())
				mockHttpClient.AssertNotCalled(ginkgo.GinkgoT(), "Do")
			})
		})

		ginkgo.When("is called with a non-pointer response object", func() {
			ginkgo.It("returns an error", func() {
				client := generic.NewClientWithHttpClient(mockHttpClient)
				_, err := client.Do(nil, 1, nil)
				gomega.Expect(err).ToNot(gomega.BeNil())
				mockHttpClient.AssertNotCalled(ginkgo.GinkgoT(), "Do")
			})
		})

		ginkgo.When("is called with a nil error response object", func() {
			ginkgo.It("returns an error", func() {
				client := generic.NewClientWithHttpClient(mockHttpClient)
				_, err := client.Do(nil, &struct{}{}, nil)
				gomega.Expect(err).ToNot(gomega.BeNil())
				mockHttpClient.AssertNotCalled(ginkgo.GinkgoT(), "Do")
			})
		})

		ginkgo.When("is called with a non-pointer error response object", func() {
			ginkgo.It("returns an error", func() {
				client := generic.NewClientWithHttpClient(mockHttpClient)
				_, err := client.Do(nil, &struct{}{}, 1)
				gomega.Expect(err).ToNot(gomega.BeNil())
				mockHttpClient.AssertNotCalled(ginkgo.GinkgoT(), "Do")
			})
		})

		ginkgo.When("is called with a nil request", func() {
			ginkgo.It("returns an error", func() {
				client := generic.NewClientWithHttpClient(mockHttpClient)
				_, err := client.Do(nil, &struct{}{}, &struct{}{})
				gomega.Expect(err).ToNot(gomega.BeNil())
				mockHttpClient.AssertNotCalled(ginkgo.GinkgoT(), "Do")
			})
		})

		ginkgo.When("is called with a valid request", func() {
			ginkgo.Context("and the request succeeds", func() {
				ginkgo.It("parses the response", func() {
					mockResponse := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`{"foo":"bar"}`))),
					}

					mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(mockResponse, nil)

					client := generic.NewClientWithHttpClient(mockHttpClient)
					var response Response
					var errorResponse Response
					code, err := client.Do(&http.Request{}, &response, &errorResponse)

					gomega.Expect(err).To(gomega.BeNil())
					gomega.Expect(code).To(gomega.Equal(http.StatusOK))
					gomega.Expect(response.Foo).To(gomega.Equal("bar"))
					gomega.Expect(errorResponse.Foo).To(gomega.Equal(""))
				})
			})
			ginkgo.Context("and the request succeeds with an un-parsable response", func() {
				ginkgo.It("returns an error", func() {
					mockResponse := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`{"foo":}`))),
					}

					mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(mockResponse, nil)

					client := generic.NewClientWithHttpClient(mockHttpClient)
					var response Response
					var errorResponse Response
					_, err := client.Do(&http.Request{}, &response, &errorResponse)

					gomega.Expect(err).ToNot(gomega.BeNil())
				})
			})
			ginkgo.Context("and the request fails", func() {
				ginkgo.It("parses the error response", func() {
					mockResponse := &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       io.NopCloser(bytes.NewReader([]byte(`{"foo":"bar"}`))),
					}

					mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(mockResponse, nil)

					client := generic.NewClientWithHttpClient(mockHttpClient)
					var response Response
					var errorResponse Response
					code, err := client.Do(&http.Request{}, &response, &errorResponse)

					gomega.Expect(err).To(gomega.BeNil())
					gomega.Expect(code).To(gomega.Equal(http.StatusInternalServerError))
					gomega.Expect(response.Foo).To(gomega.Equal(""))
					gomega.Expect(errorResponse.Foo).To(gomega.Equal("bar"))
				})
			})
			ginkgo.Context("and the request fails with an un-parsable error response", func() {
				ginkgo.It("returns an error", func() {
					mockResponse := &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       io.NopCloser(bytes.NewReader([]byte(`{"foo":}`))),
					}

					mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(mockResponse, nil)

					client := generic.NewClientWithHttpClient(mockHttpClient)
					var response Response
					var errorResponse Response
					_, err := client.Do(&http.Request{}, &response, &errorResponse)

					gomega.Expect(err).ToNot(gomega.BeNil())
				})
			})
		})
	})
})
