package tournament_test

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
	"pinman/internal/app/api/tournament"
	"pinman/internal/app/generated"
	"pinman/internal/models"
	"pinman/internal/utils"
	"regexp"
	"testing"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestTournament(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Tournament Suite")
}

var _ = ginkgo.Describe("Controller", func() {
	var controller *tournament.Controller
	var db *gorm.DB
	var mock sqlmock.Sqlmock
	var rr *httptest.ResponseRecorder
	var router *gin.Engine
	var _ *gin.Context
	var userObj *models.User

	ginkgo.BeforeEach(func() {
		db, mock = utils.NewGormMock()
		controller = tournament.NewController(db)
		_, rr, router = utils.NewGinTestCtx()

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
	})

	ginkgo.When("CreateTournament is called", func() {
		var payload *generated.TournamentCreate
		ginkgo.BeforeEach(func() {
			settings := generated.MultiRoundTournamentSettings{
				GamesPerRound:       4,
				LowestScoresDropped: 3,
				Rounds:              8,
			}
			si := &generated.TournamentSettings{}
			err := si.FromMultiRoundTournamentSettings(settings)
			gomega.Expect(err).To(gomega.BeNil())

			payload = &generated.TournamentCreate{
				LeagueId:   uuid.New().String(),
				LocationId: uuid.New().String(),
				Name:       "Test Tournament",
				Slug:       "test-tournament",
				Settings:   *si,
				Type:       generated.MultiRoundTournament,
			}
		})

		ginkgo.Context("with a valid payload", func() {
			ginkgo.It("returns a 201", func() {
				router.Use(func(c *gin.Context) {
					c.Set("user", userObj)
				})

				const leaguesQuery = `SELECT * FROM "leagues" WHERE id = $1 ORDER BY "leagues"."id" LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(leaguesQuery)).
					WithArgs(payload.LeagueId).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(payload.LeagueId),
					)

				const locationsQuery = `SELECT * FROM "locations" WHERE id = $1 ORDER BY "locations"."id" LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(locationsQuery)).
					WithArgs(payload.LocationId).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(payload.LocationId),
					)

				insertedSettings, err := payload.Settings.MarshalJSON()
				gomega.Expect(err).To(gomega.BeNil())
				mock.ExpectBegin()
				const sqlInsert = `INSERT INTO "tournaments" ("name","slug","type","settings","location_id","league_id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
				mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
					WithArgs(payload.Name, payload.Slug, payload.Type, string(insertedSettings), payload.LocationId, payload.LeagueId).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(uuid.New()),
					)
				mock.ExpectCommit()

				body, err := json.Marshal(payload)
				gomega.Expect(err).To(gomega.BeNil())
				req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
				gomega.Expect(err).To(gomega.BeNil())

				router.POST("/", controller.CreateTournament)
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusCreated))
				response := &generated.TournamentResponse{}
				err = json.Unmarshal(rr.Body.Bytes(), response)
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(response.Tournament.Name).To(gomega.Equal(payload.Name))
				gomega.Expect(response.Tournament.Slug).To(gomega.Equal(payload.Slug))
				gomega.Expect(response.Tournament.Type).To(gomega.Equal(payload.Type))
			})
		})
		ginkgo.Context("with a location id that does not exist", func() {
			ginkgo.It("returns a 400", func() {
				router.Use(func(c *gin.Context) {
					c.Set("user", userObj)
				})

				const leaguesQuery = `SELECT * FROM "leagues" WHERE id = $1 ORDER BY "leagues"."id" LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(leaguesQuery)).
					WithArgs(payload.LeagueId).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(payload.LeagueId),
					)

				const locationsQuery = `SELECT * FROM "locations" WHERE id = $1 ORDER BY "locations"."id" LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(locationsQuery)).
					WithArgs(payload.LocationId).
					WillReturnError(gorm.ErrRecordNotFound)

				body, err := json.Marshal(payload)
				gomega.Expect(err).To(gomega.BeNil())
				req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
				gomega.Expect(err).To(gomega.BeNil())

				router.POST("/", controller.CreateTournament)
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusBadRequest))
			})
		})
		ginkgo.Context("with a league id that does not exist", func() {
			ginkgo.It("returns a 400", func() {
				router.Use(func(c *gin.Context) {
					c.Set("user", userObj)
				})

				const leaguesQuery = `SELECT * FROM "leagues" WHERE id = $1 ORDER BY "leagues"."id" LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(leaguesQuery)).
					WithArgs(payload.LeagueId).
					WillReturnError(gorm.ErrRecordNotFound)

				const locationsQuery = `SELECT * FROM "locations" WHERE id = $1 ORDER BY "locations"."id" LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(locationsQuery)).
					WithArgs(payload.LocationId).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(payload.LocationId),
					)

				body, err := json.Marshal(payload)
				gomega.Expect(err).To(gomega.BeNil())
				req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
				gomega.Expect(err).To(gomega.BeNil())

				router.POST("/", controller.CreateTournament)
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusBadRequest))
			})
		})
		ginkgo.Context("with a payload that is not valid", func() {
			ginkgo.It("returns a 400", func() {
				router.Use(func(c *gin.Context) {
					c.Set("user", userObj)
				})

				body, err := json.Marshal(struct{}{})
				gomega.Expect(err).To(gomega.BeNil())
				req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
				gomega.Expect(err).To(gomega.BeNil())

				router.POST("/", controller.CreateTournament)
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusBadRequest))
			})
		})
		ginkgo.Context("with a settings payload that is invalid", func() {
			ginkgo.It("returns a 400", func() {
				router.Use(func(c *gin.Context) {
					c.Set("user", userObj)
				})

				tpl := `{
"name":"test","slug":"test","type":"%s",
"league_id":"%s","location_id":"%s","settings":{"test":"test"}
}`
				body := []byte(fmt.Sprintf(tpl, payload.Type, payload.LeagueId, payload.LocationId))

				req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
				gomega.Expect(err).To(gomega.BeNil())

				router.POST("/", controller.CreateTournament)
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusBadRequest))
				gomega.Expect(rr.Body.String()).To(gomega.ContainSubstring("TournamentSettings"))
			})
		})
		ginkgo.Context("with existing slug", func() {
			ginkgo.It("returns a 400", func() {
				router.Use(func(c *gin.Context) {
					c.Set("user", userObj)
				})

				const leaguesQuery = `SELECT * FROM "leagues" WHERE id = $1 ORDER BY "leagues"."id" LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(leaguesQuery)).
					WithArgs(payload.LeagueId).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(payload.LeagueId),
					)

				const locationsQuery = `SELECT * FROM "locations" WHERE id = $1 ORDER BY "locations"."id" LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(locationsQuery)).
					WithArgs(payload.LocationId).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(payload.LocationId),
					)

				insertedSettings, err := payload.Settings.MarshalJSON()
				gomega.Expect(err).To(gomega.BeNil())

				mock.ExpectBegin()
				const sqlInsert = `INSERT INTO "tournaments" ("name","slug","type","settings","location_id","league_id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
				mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
					WithArgs(payload.Name, payload.Slug, payload.Type, string(insertedSettings), payload.LocationId, payload.LeagueId).
					WillReturnError(fmt.Errorf("ERROR: duplicate key value violates unique constraint \"idx_tournament_slug\" (SQLSTATE 23505)"))
				mock.ExpectRollback()

				body, err := json.Marshal(payload)
				gomega.Expect(err).To(gomega.BeNil())
				req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
				gomega.Expect(err).To(gomega.BeNil())

				router.POST("/", controller.CreateTournament)
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusBadRequest))
			})
		})
		ginkgo.Context("with a database error", func() {
			ginkgo.It("returns a 500", func() {
				router.Use(func(c *gin.Context) {
					c.Set("user", userObj)
				})

				const leaguesQuery = `SELECT * FROM "leagues" WHERE id = $1 ORDER BY "leagues"."id" LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(leaguesQuery)).
					WithArgs(payload.LeagueId).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(payload.LeagueId),
					)

				const locationsQuery = `SELECT * FROM "locations" WHERE id = $1 ORDER BY "locations"."id" LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(locationsQuery)).
					WithArgs(payload.LocationId).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(payload.LocationId),
					)

				insertedSettings, err := payload.Settings.MarshalJSON()
				gomega.Expect(err).To(gomega.BeNil())

				mock.ExpectBegin()
				const sqlInsert = `INSERT INTO "tournaments" ("name","slug","type","settings","location_id","league_id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
				mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
					WithArgs(payload.Name, payload.Slug, payload.Type, string(insertedSettings), payload.LocationId, payload.LeagueId).
					WillReturnError(fmt.Errorf("ERROR: database error"))
				mock.ExpectRollback()

				body, err := json.Marshal(payload)
				gomega.Expect(err).To(gomega.BeNil())
				req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
				gomega.Expect(err).To(gomega.BeNil())

				router.POST("/", controller.CreateTournament)
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusInternalServerError))
			})
		})
	})
	ginkgo.Describe("ListTournaments", func() {
		ginkgo.Context("with a valid request", func() {
			ginkgo.It("returns a 200", func() {
				router.Use(func(c *gin.Context) {
					c.Set("user", userObj)
				})

				mockMultiRoundSettings := generated.MultiRoundTournamentSettings{
					GamesPerRound:       4,
					Rounds:              8,
					LowestScoresDropped: 3,
				}
				mockSettings, err := json.Marshal(mockMultiRoundSettings)
				gomega.Expect(err).To(gomega.BeNil())

				mockTournament := models.Tournament{
					ID:         uuid.New(),
					Name:       "Test Tournament",
					Slug:       "test-tournament",
					Type:       generated.MultiRoundTournament,
					Settings:   mockSettings,
					LocationID: uuid.New(),
					LeagueID:   uuid.New(),
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}

				const query = `SELECT * FROM "tournaments"`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(
						sqlmock.NewRows([]string{
							"id", "name", "slug", "type", "settings", "location_id", "league_id", "created_at", "updated_at",
						}).
							AddRow(
								mockTournament.ID.String(), mockTournament.Name, mockTournament.Slug,
								mockTournament.Type, mockTournament.Settings, mockTournament.LocationID.String(),
								mockTournament.LeagueID.String(), mockTournament.CreatedAt, mockTournament.UpdatedAt,
							),
					)

				const leaguesQuery = `SELECT * FROM "leagues"`
				mock.ExpectQuery(regexp.QuoteMeta(leaguesQuery)).
					WithArgs(mockTournament.LeagueID.String()).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(mockTournament.LeagueID.String()),
					)

				const locationsQuery = `SELECT * FROM "locations"`
				mock.ExpectQuery(regexp.QuoteMeta(locationsQuery)).
					WithArgs(mockTournament.LocationID.String()).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(mockTournament.LocationID.String()),
					)

				req, err := http.NewRequest(http.MethodGet, "/", nil)
				gomega.Expect(err).To(gomega.BeNil())

				router.GET("/", controller.ListTournaments)
				router.ServeHTTP(rr, req)

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
				gomega.Expect(rr.Body.String()).To(gomega.ContainSubstring(mockTournament.Name))
			})
		})
	})
})
