package league_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"pinman/internal/app/api/auth"
	"pinman/internal/app/api/league"
	"pinman/internal/app/generated"
	"pinman/internal/models"
	"pinman/internal/utils"
	"regexp"
	"testing"
	"time"
)

func TestLeague(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "League Suite")
}

var _ = ginkgo.Describe("Controller", func() {
	var _ *gin.Context
	var rr *httptest.ResponseRecorder
	var router *gin.Engine
	var db *gorm.DB
	var mock sqlmock.Sqlmock
	var controller *league.Controller
	var userObj *models.User
	var locationObj *models.Location

	ginkgo.BeforeEach(func() {
		_, rr, router = utils.NewGinTestCtx()
		db, mock = utils.NewGormMock()
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

		locationObj = &models.Location{
			ID:           uuid.New(),
			Name:         "Test Location",
			Slug:         "test-location",
			Address:      "123 Test St",
			PinballMapID: 123,
			CreatedAt:    time.Now().Add(-1 * time.Hour),
			UpdatedAt:    time.Now(),
		}

		controller = league.NewController(db)
	})

	ginkgo.When("CreateLeague receives a request", func() {
		var payload *generated.LeagueCreate
		ginkgo.BeforeEach(func() {
			payload = &generated.LeagueCreate{
				Name:       "Test League",
				Slug:       "test-league",
				LocationId: locationObj.ID.String(),
			}
		})

		ginkgo.Context("with valid payload", func() {
			ginkgo.It("succeeds", func() {
				router.Use(func(ctx *gin.Context) {
					ctx.Set(auth.IdentityKey, userObj)
				})

				const sqlQuery = `SELECT * FROM "locations" WHERE id = $1 ORDER BY "locations"."id" LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(payload.LocationId).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "slug", "address", "pinball_map_id", "created_at", "updated_at"}).
							AddRow(locationObj.ID.String(), locationObj.Name, locationObj.Slug, locationObj.Address, locationObj.PinballMapID, locationObj.CreatedAt, locationObj.UpdatedAt),
					)

				mock.ExpectBegin()
				const sqlInsert = `INSERT INTO "leagues" ("name","slug","owner_id","location_id","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
				mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
					WithArgs(payload.Name, payload.Slug, userObj.ID.String(), payload.LocationId, utils.AnyTime{}, utils.AnyTime{}).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(uuid.New()),
					)
				mock.ExpectCommit()

				body, err := json.Marshal(payload)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				req, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))

				router.POST("/", controller.CreateLeague)
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusCreated))
				response := &generated.LeagueResponse{}
				err = json.Unmarshal(rr.Body.Bytes(), response)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(mock.ExpectationsWereMet()).ToNot(gomega.HaveOccurred())

				gomega.Expect(response.League.Name).To(gomega.Equal(payload.Name))
				gomega.Expect(response.League.Slug).To(gomega.Equal(payload.Slug))
				gomega.Expect(response.League.LocationId).To(gomega.Equal(payload.LocationId))

			})
		})

		ginkgo.Context("with a location id that does not exist", func() {
			ginkgo.It("fails with 400 bad request", func() {
				router.Use(func(ctx *gin.Context) {
					ctx.Set(auth.IdentityKey, userObj)
				})

				const sqlQuery = `SELECT * FROM "locations" WHERE id = $1 ORDER BY "locations"."id" LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(payload.LocationId).
					WillReturnError(gorm.ErrRecordNotFound)

				body, err := json.Marshal(payload)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				req, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))

				router.POST("/", controller.CreateLeague)
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusBadRequest))
				response := &generated.ErrorResponse{}
				err = json.Unmarshal(rr.Body.Bytes(), response)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(mock.ExpectationsWereMet()).ToNot(gomega.HaveOccurred())
				gomega.Expect(response.Detail).To(gomega.ContainSubstring("does not exist"))
			})
		})

		ginkgo.Context("when the location lookup query fails", func() {
			ginkgo.It("fails with 500", func() {
				router.Use(func(ctx *gin.Context) {
					ctx.Set(auth.IdentityKey, userObj)
				})

				const sqlQuery = `SELECT * FROM "locations" WHERE id = $1 ORDER BY "locations"."id" LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(payload.LocationId).
					WillReturnError(errors.New("some error"))

				body, err := json.Marshal(payload)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				req, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))

				router.POST("/", controller.CreateLeague)
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusInternalServerError))
				response := &generated.ErrorResponse{}
				err = json.Unmarshal(rr.Body.Bytes(), response)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(mock.ExpectationsWereMet()).ToNot(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("with invalid payload", func() {
			ginkgo.It("fails with bad request", func() {
				router.Use(func(ctx *gin.Context) {
					ctx.Set(auth.IdentityKey, userObj)
				})

				payload := map[string]interface{}{
					"foo": "bar",
					"int": 64,
				}

				body, err := json.Marshal(payload)
				gomega.Expect(err).To(gomega.BeNil())
				req, err := http.NewRequest("POST", "/", bytes.NewReader(body))

				router.POST("/", controller.CreateLeague)
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusBadRequest))

				response := &generated.BadRequest{}
				err = json.Unmarshal(rr.Body.Bytes(), response)
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with existing slug", func() {
			ginkgo.It("fails with conflict", func() {
				router.Use(func(ctx *gin.Context) {
					ctx.Set(auth.IdentityKey, userObj)
				})

				const sqlQuery = `SELECT * FROM "locations" WHERE id = $1 ORDER BY "locations"."id" LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(payload.LocationId).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(payload.LocationId),
					)

				mock.ExpectBegin()
				const sqlInsert = `INSERT INTO "leagues" ("name","slug","owner_id","location_id","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
				mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
					WithArgs(payload.Name, payload.Slug, userObj.ID.String(), payload.LocationId, utils.AnyTime{}, utils.AnyTime{}).
					WillReturnError(fmt.Errorf("duplicate key value violates unique"))
				mock.ExpectRollback()

				body, err := json.Marshal(payload)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				req, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))

				router.POST("/", controller.CreateLeague)
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusConflict))
				response := &generated.ErrorResponse{}
				err = json.Unmarshal(rr.Body.Bytes(), response)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(mock.ExpectationsWereMet()).ToNot(gomega.HaveOccurred())

				gomega.Expect(response.Detail).To(gomega.Equal("league with slug already exists"))
			})
		})

		ginkgo.Context("with unknown sql error", func() {
			ginkgo.It("fails", func() {
				router.Use(func(ctx *gin.Context) {
					ctx.Set(auth.IdentityKey, userObj)
				})

				const sqlQuery = `SELECT * FROM "locations" WHERE id = $1 ORDER BY "locations"."id" LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(payload.LocationId).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "slug", "address", "pinball_map_id", "created_at", "updated_at"}).
							AddRow(locationObj.ID.String(), locationObj.Name, locationObj.Slug, locationObj.Address, locationObj.PinballMapID, locationObj.CreatedAt, locationObj.UpdatedAt),
					)

				mock.ExpectBegin()
				const sqlInsert = `INSERT INTO "leagues" ("name","slug","owner_id","location_id","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
				mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
					WithArgs(payload.Name, payload.Slug, userObj.ID.String(), payload.LocationId, utils.AnyTime{}, utils.AnyTime{}).
					WillReturnError(fmt.Errorf("unknown error"))
				mock.ExpectRollback()

				body, err := json.Marshal(payload)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				req, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))

				router.POST("/", controller.CreateLeague)
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusInternalServerError))
				response := &generated.ErrorResponse{}
				err = json.Unmarshal(rr.Body.Bytes(), response)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(mock.ExpectationsWereMet()).ToNot(gomega.HaveOccurred())

				gomega.Expect(response.Detail).To(gomega.Equal("failed to create league"))
			})
		})
	})

	ginkgo.When("ListLeagues receives a request", func() {
		ginkgo.Context("with valid request", func() {
			ginkgo.It("succeeds", func() {
				router.Use(func(ctx *gin.Context) {
					ctx.Set(auth.IdentityKey, userObj)
				})

				leagueObj := &models.League{
					ID:         uuid.New(),
					Name:       "Test League",
					Slug:       "test-league",
					Owner:      *userObj,
					LocationID: uuid.New(),
					CreatedAt:  time.Now().Add(-1 * time.Hour),
					UpdatedAt:  time.Now(),
				}

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "leagues"`)).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "slug", "owner_id", "location_id", "created_at", "updated_at"}).
							AddRow(leagueObj.ID.String(), leagueObj.Name, leagueObj.Slug, leagueObj.Owner.ID, leagueObj.LocationID, leagueObj.CreatedAt, leagueObj.UpdatedAt),
					)

				req, err := http.NewRequest("GET", "/", nil)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				router.GET("/", controller.ListLeagues)
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
				response := &generated.LeagueListResponse{}
				err = json.Unmarshal(rr.Body.Bytes(), response)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(mock.ExpectationsWereMet()).ToNot(gomega.HaveOccurred())

				gomega.Expect(response.Leagues).To(gomega.HaveLen(1))
				gomega.Expect(response.Leagues[0].Name).To(gomega.Equal(leagueObj.Name))
				gomega.Expect(response.Leagues[0].Slug).To(gomega.Equal(leagueObj.Slug))
				gomega.Expect(response.Leagues[0].LocationId).To(gomega.Equal(leagueObj.LocationID.String()))
			})
		})
		ginkgo.Context("with unknown sql error", func() {
			ginkgo.It("fails", func() {
				router.Use(func(ctx *gin.Context) {
					ctx.Set(auth.IdentityKey, userObj)
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "leagues"`)).
					WillReturnError(fmt.Errorf("unknown error"))

				req, err := http.NewRequest("GET", "/", nil)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				router.GET("/", controller.ListLeagues)
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusInternalServerError))
				response := &generated.ErrorResponse{}
				err = json.Unmarshal(rr.Body.Bytes(), response)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(mock.ExpectationsWereMet()).ToNot(gomega.HaveOccurred())

				gomega.Expect(response.Detail).To(gomega.Equal("failed to list leagues"))
			})
		})
	})

	ginkgo.When("GetLeagueWithSlug receives a request", func() {
		ginkgo.Context("with valid payload", func() {
			ginkgo.It("succeeds", func() {
				leagueObj := &models.League{
					ID:         uuid.New(),
					Name:       "Test League",
					Slug:       "test-league",
					Owner:      *userObj,
					LocationID: locationObj.ID,
					CreatedAt:  time.Now().Add(-1 * time.Hour),
					UpdatedAt:  time.Now(),
				}

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "leagues" WHERE slug = $1`)).
					WithArgs(leagueObj.Slug).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "slug", "owner_id", "location_id", "created_at", "updated_at"}).
							AddRow(leagueObj.ID.String(), leagueObj.Name, leagueObj.Slug, leagueObj.Owner.ID.String(), leagueObj.LocationID.String(), leagueObj.CreatedAt, leagueObj.UpdatedAt),
					)

				req, err := http.NewRequest("GET", "/", nil)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				router.GET("/", func(c *gin.Context) {
					controller.GetLeagueWithSlug(c, leagueObj.Slug)
				})
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
				response := &generated.LeagueResponse{}
				err = json.Unmarshal(rr.Body.Bytes(), response)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(mock.ExpectationsWereMet()).ToNot(gomega.HaveOccurred())

				gomega.Expect(response.League.Name).To(gomega.Equal(leagueObj.Name))
				gomega.Expect(response.League.Slug).To(gomega.Equal(leagueObj.Slug))
				gomega.Expect(response.League.LocationId).To(gomega.Equal(leagueObj.LocationID.String()))
			})
		})
		ginkgo.Context("when location is not found", func() {
			ginkgo.It("fails", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "leagues" WHERE slug = $1`)).
					WillReturnError(gorm.ErrRecordNotFound)

				req, err := http.NewRequest("GET", "/", nil)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				router.GET("/", func(c *gin.Context) {
					controller.GetLeagueWithSlug(c, "foo-bar")
				})
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusNotFound))
				response := &generated.ErrorResponse{}
				err = json.Unmarshal(rr.Body.Bytes(), response)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(mock.ExpectationsWereMet()).ToNot(gomega.HaveOccurred())

				gomega.Expect(response.Detail).To(gomega.Equal("league not found"))

			})
		})
		ginkgo.Context("with unknown sql error", func() {
			ginkgo.It("fails", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "leagues" WHERE slug = $1`)).
					WillReturnError(fmt.Errorf("unknown error"))

				req, err := http.NewRequest("GET", "/", nil)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				router.GET("/", func(c *gin.Context) {
					controller.GetLeagueWithSlug(c, "foo-bar")
				})
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusInternalServerError))
				response := &generated.ErrorResponse{}
				err = json.Unmarshal(rr.Body.Bytes(), response)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(mock.ExpectationsWereMet()).ToNot(gomega.HaveOccurred())

				gomega.Expect(response.Detail).To(gomega.Equal("failed to get league"))

			})
		})
	})
})
