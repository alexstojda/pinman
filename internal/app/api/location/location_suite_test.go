package location_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"pinman/internal/app/api/auth"
	"pinman/internal/app/api/location"
	"pinman/internal/app/generated"
	"pinman/internal/clients/pinballmap"
	"pinman/internal/models"
	"pinman/internal/utils"
	"regexp"
	"testing"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestLocation(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Location Suite")
}

var _ = ginkgo.Describe("NewController", func() {
	ginkgo.It("should return a new controller", func() {
		db, _ := utils.NewGormMock()
		controller := location.NewController(db)
		gomega.Expect(controller).ToNot(gomega.BeNil())
		gomega.Expect(controller.DB).ToNot(gomega.BeNil())
	})
})

var _ = ginkgo.Describe("Location", func() {
	var controller *location.Controller
	var db *gorm.DB
	var mock sqlmock.Sqlmock
	var rr *httptest.ResponseRecorder
	var router *gin.Engine
	var ctx *gin.Context
	var userObj *models.User
	var mockPinballMapClient *pinballmap.MockClientInterface
	var mockPinballLocationsResponse *pinballmap.Location

	ginkgo.BeforeEach(func() {
		db, mock = utils.NewGormMock()
		mockPinballMapClient = pinballmap.NewMockClientInterface(ginkgo.GinkgoT())
		controller = location.NewControllerWithClient(db, mockPinballMapClient)
		ctx, rr, router = utils.NewGinTestCtx()

		userObj = &models.User{
			ID:        uuid.New(),
			Name:      "John Doe",
			Email:     "email@example.com",
			Password:  "password",
			Role:      "user",
			Verified:  false,
			CreatedAt: time.Now().Add(-1 * time.Hour),
			UpdatedAt: time.Now(),
		}

		mockPinballLocationsResponse = &pinballmap.Location{
			ID:      1,
			Name:    "Pinballz Arcade",
			Street:  "123 Main St",
			City:    "Austin",
			State:   "TX",
			Country: "USA",
		}
	})

	ginkgo.Describe("CreateLocation", func() {
		ginkgo.Context("is called with a valid payload", func() {
			ginkgo.It("returns a 201", func() {
				router.Use(func(ctx *gin.Context) {
					ctx.Set(auth.IdentityKey, userObj)
				})
				mockPinballMapClient.On("GetLocation", 1).Return(mockPinballLocationsResponse, nil)

				mock.ExpectBegin()
				const sqlInsert = `INSERT INTO "locations" ("name","slug","address","pinball_map_id","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
				mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
					WithArgs(
						"Pinballz Arcade",
						"pinballz-arcade",
						"123 Main St, Austin, TX, USA",
						1,
						utils.AnyTime{}, utils.AnyTime{},
					).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
				mock.ExpectCommit()

				body, err := json.Marshal(&generated.LocationCreate{
					PinballMapId: mockPinballLocationsResponse.ID,
				})
				gomega.Expect(err).To(gomega.BeNil())
				req, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
				gomega.Expect(err).To(gomega.BeNil())

				router.POST("/", controller.CreateLocation)
				router.ServeHTTP(rr, req)
				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusCreated))
				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())

				response := &generated.LocationResponse{}
				err = json.Unmarshal(rr.Body.Bytes(), response)
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(response.Location).ToNot(gomega.BeNil())

				gomega.Expect(response.Location.Name).To(gomega.Equal(mockPinballLocationsResponse.Name))
			})
		})

		ginkgo.Context("is called with an invalid payload", func() {
			ginkgo.It("returns a 400", func() {
				router.Use(func(ctx *gin.Context) {
					ctx.Set(auth.IdentityKey, userObj)
				})

				req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte("invalid json")))
				gomega.Expect(err).To(gomega.BeNil())

				router.POST("/", controller.CreateLocation)
				router.ServeHTTP(rr, req)
				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusBadRequest))
			})
		})

		ginkgo.Context("is called with a valid payload but the pinball map client returns an error", func() {
			ginkgo.It("returns a 500", func() {
				router.Use(func(ctx *gin.Context) {
					ctx.Set(auth.IdentityKey, userObj)
				})
				mockPinballMapClient.On("GetLocation", 1).Return(nil, fmt.Errorf("some error"))

				body, err := json.Marshal(&generated.LocationCreate{
					PinballMapId: 1,
				})
				gomega.Expect(err).To(gomega.BeNil())
				req, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
				gomega.Expect(err).To(gomega.BeNil())

				router.POST("/", controller.CreateLocation)
				router.ServeHTTP(rr, req)
				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusInternalServerError))
			})
		})

		ginkgo.Context("is called with a location whose slug already exists", func() {
			ginkgo.It("returns a 409", func() {
				router.Use(func(ctx *gin.Context) {
					ctx.Set(auth.IdentityKey, userObj)
				})
				mockPinballMapClient.On("GetLocation", 1).Return(mockPinballLocationsResponse, nil)

				mock.ExpectBegin()
				const sqlInsert = `INSERT INTO "locations" ("name","slug","address","pinball_map_id","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
				mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
					WithArgs(
						"Pinballz Arcade",
						"pinballz-arcade",
						"123 Main St, Austin, TX, USA",
						1,
						utils.AnyTime{}, utils.AnyTime{},
					).WillReturnError(fmt.Errorf("duplicate key value violates unique"))
				mock.ExpectRollback()

				body, err := json.Marshal(&generated.LocationCreate{
					PinballMapId: mockPinballLocationsResponse.ID,
				})
				gomega.Expect(err).To(gomega.BeNil())
				req, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
				gomega.Expect(err).To(gomega.BeNil())

				router.POST("/", controller.CreateLocation)
				router.ServeHTTP(rr, req)
				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusConflict))
				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())
			})
		})

		ginkgo.Context("and the query fails", func() {
			ginkgo.It("returns a 500", func() {
				router.Use(func(ctx *gin.Context) {
					ctx.Set(auth.IdentityKey, userObj)
				})
				mockPinballMapClient.On("GetLocation", 1).Return(mockPinballLocationsResponse, nil)

				mock.ExpectBegin()
				const sqlInsert = `INSERT INTO "locations" ("name","slug","address","pinball_map_id","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
				mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
					WithArgs(
						"Pinballz Arcade",
						"pinballz-arcade",
						"123 Main St, Austin, TX, USA",
						1,
						utils.AnyTime{}, utils.AnyTime{},
					).WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()

				body, err := json.Marshal(&generated.LocationCreate{
					PinballMapId: mockPinballLocationsResponse.ID,
				})
				gomega.Expect(err).To(gomega.BeNil())
				req, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
				gomega.Expect(err).To(gomega.BeNil())

				router.POST("/", controller.CreateLocation)
				router.ServeHTTP(rr, req)
				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusInternalServerError))
				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("GetLocation", func() {
		ginkgo.Context("is called with a valid location slug", func() {
			ginkgo.It("returns a 200", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "locations" WHERE slug = $1`)).
					WithArgs("pinballz-arcade").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "slug", "address", "pinball_map_id", "created_at", "updated_at"}).
						AddRow(uuid.New(), "Pinballz Arcade", "pinballz-arcade", "123 Main St, Austin, TX, USA", 1, time.Now(), time.Now()))

				controller.GetLocationWithSlug(ctx, "pinballz-arcade")

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())

				response := &generated.LocationResponse{}
				err := json.Unmarshal(rr.Body.Bytes(), response)
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(response.Location).ToNot(gomega.BeNil())

				gomega.Expect(response.Location.Name).To(gomega.Equal("Pinballz Arcade"))
			})
		})
		ginkgo.Context("is called with an unknown location slug", func() {
			ginkgo.It("returns a 404", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "locations" WHERE slug = $1`)).
					WithArgs("pinballz-arcade").
					WillReturnError(gorm.ErrRecordNotFound)

				controller.GetLocationWithSlug(ctx, "pinballz-arcade")

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusNotFound))
				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())
			})
		})
		ginkgo.Context("and the query fails", func() {
			ginkgo.It("returns a 500", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "locations" WHERE slug = $1`)).
					WithArgs("pinballz-arcade").
					WillReturnError(fmt.Errorf("some error"))

				controller.GetLocationWithSlug(ctx, "pinballz-arcade")

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusInternalServerError))
				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("ListLocations", func() {
		ginkgo.Context("is called", func() {
			ginkgo.It("returns a 200", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "locations"`)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "slug", "address", "pinball_map_id", "created_at", "updated_at"}).
						AddRow(uuid.New(), "Pinballz Arcade", "pinballz-arcade", "123 Main St, Austin, TX, USA", 1, time.Now(), time.Now()))

				controller.ListLocations(ctx)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())

				response := &generated.LocationListResponse{}
				err := json.Unmarshal(rr.Body.Bytes(), response)
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(response.Locations).ToNot(gomega.BeNil())

				gomega.Expect(response.Locations[0].Name).To(gomega.Equal("Pinballz Arcade"))
			})
		})
		ginkgo.Context("and the query fails", func() {
			ginkgo.It("returns a 500", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "locations"`)).
					WillReturnError(fmt.Errorf("some error"))

				controller.ListLocations(ctx)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusInternalServerError))
				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())
			})
		})
	})
})
