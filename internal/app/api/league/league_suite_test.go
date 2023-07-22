package league_test

import (
	"bytes"
	"encoding/json"
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

		controller = league.NewController(db)
	})

	ginkgo.When("CreateLeague receives a request", func() {
		payload := &generated.LeagueCreate{
			Name:     "Test League",
			Slug:     "test-league",
			Location: "Test Location",
		}

		ginkgo.Context("with valid payload", func() {
			ginkgo.It("succeeds", func() {
				router.Use(func(ctx *gin.Context) {
					ctx.Set(auth.IdentityKey, userObj)
				})

				mock.ExpectBegin()
				const sqlInsert = `INSERT INTO "leagues" ("name","slug","owner_id","location","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
				mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
					WithArgs(payload.Name, payload.Slug, userObj.ID.String(), payload.Location, utils.AnyTime{}, utils.AnyTime{}).
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
				gomega.Expect(response.League.Location).To(gomega.Equal(payload.Location))

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

				mock.ExpectBegin()
				const sqlInsert = `INSERT INTO "leagues" ("name","slug","owner_id","location","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
				mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
					WithArgs(payload.Name, payload.Slug, userObj.ID.String(), payload.Location, utils.AnyTime{}, utils.AnyTime{}).
					WillReturnError(fmt.Errorf("duplicate key value violates unique"))

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

				mock.ExpectBegin()
				const sqlInsert = `INSERT INTO "leagues" ("name","slug","owner_id","location","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
				mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
					WithArgs(payload.Name, payload.Slug, userObj.ID.String(), payload.Location, utils.AnyTime{}, utils.AnyTime{}).
					WillReturnError(fmt.Errorf("unknown error"))

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
})
