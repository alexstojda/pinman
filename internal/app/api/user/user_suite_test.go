package user_test

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
	"pinman/internal/app/api/user"
	"pinman/internal/app/generated"
	"pinman/internal/models"
	"pinman/internal/utils"
	"regexp"
	"testing"
	"time"

	g "github.com/onsi/ginkgo/v2"
	m "github.com/onsi/gomega"
)

func TestUser(t *testing.T) {
	m.RegisterFailHandler(g.Fail)
	g.RunSpecs(t, "User Suite")
}

var _ = g.Describe("Controller", func() {
	var ctx *gin.Context
	var rr *httptest.ResponseRecorder
	var router *gin.Engine
	var db *gorm.DB
	var mock sqlmock.Sqlmock
	var controller *user.Controller

	g.BeforeEach(func() {
		ctx, rr, router = utils.NewGinTestCtx()
		db, mock = utils.NewGormMock()

		controller = user.NewController(db)
	})

	g.When("GetMe receives a request", func() {
		userObj := &models.User{
			ID:        uuid.New(),
			Name:      "John Doe",
			Email:     "email@example.com",
			Password:  "password",
			Role:      "user",
			Verified:  false,
			CreatedAt: time.Now().Add(-1 * time.Hour),
			UpdatedAt: time.Now(),
		}

		g.It("returns authenticated user's information", func() {

			ctx.Set(auth.IdentityKey, userObj)
			controller.GetMe(ctx)

			m.Expect(rr.Code).To(m.BeEquivalentTo(http.StatusOK))

			resp := &generated.UserResponse{}
			err := json.Unmarshal(rr.Body.Bytes(), resp)
			m.Expect(err).To(m.BeNil())
			m.Expect(resp).To(m.BeEquivalentTo(
				&generated.UserResponse{
					User: &generated.User{
						Email:     userObj.Email,
						Id:        userObj.ID.String(),
						Name:      userObj.Name,
						Role:      userObj.Role,
						CreatedAt: utils.FormatTime(userObj.CreatedAt),
						UpdatedAt: utils.FormatTime(userObj.UpdatedAt),
					},
				},
			))
		})
	})

	g.When("SignUpUser receives a request", func() {
		g.Context("that is valid", func() {
			g.It("succeeds", func() {
				payload := &generated.UserRegister{
					Email:           "user@example.com",
					Name:            "Jason Doe",
					Password:        "password",
					PasswordConfirm: "password",
				}

				mock.ExpectBegin()
				const sqlInsert = `INSERT INTO "users" ("name","email","password","role","verified","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`
				mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
					WithArgs(payload.Name, payload.Email, utils.AnyString{}, "user", true, utils.AnyTime{}, utils.AnyTime{}).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(uuid.New()),
					)
				mock.ExpectCommit()

				body, err := json.Marshal(payload)
				m.Expect(err).To(m.BeNil())
				req, err := http.NewRequest("POST", "/", bytes.NewReader(body))

				router.POST("/", controller.SignUpUser)
				router.ServeHTTP(rr, req)

				m.Expect(rr.Code).To(m.Equal(http.StatusCreated))

				response := &generated.UserResponse{}
				err = json.Unmarshal(rr.Body.Bytes(), response)
				m.Expect(err).To(m.BeNil())
				m.Expect(mock.ExpectationsWereMet()).To(m.BeNil())

				m.Expect(response.User.Email).To(m.Equal(payload.Email))
				m.Expect(response.User.Name).To(m.Equal(payload.Name))
			})
		})
		g.Context("with invalid payload", func() {
			g.It("fails with bad request", func() {
				payload := map[string]interface{}{
					"foo": "bar",
					"int": 64,
				}

				body, err := json.Marshal(payload)
				m.Expect(err).To(m.BeNil())
				req, err := http.NewRequest("POST", "/", bytes.NewReader(body))

				router.POST("/", controller.SignUpUser)
				router.ServeHTTP(rr, req)

				m.Expect(rr.Code).To(m.Equal(http.StatusBadRequest))

				response := &generated.BadRequest{}
				err = json.Unmarshal(rr.Body.Bytes(), response)
				m.Expect(err).To(m.BeNil())
				m.Expect(mock.ExpectationsWereMet()).To(m.BeNil())
			})
		})
		g.Context("with invalid password confirmation", func() {
			g.It("fails with bad request", func() {
				payload := &generated.UserRegister{
					Email:           "user@example.com",
					Name:            "Jason Doe",
					Password:        "password",
					PasswordConfirm: "different password",
				}

				body, err := json.Marshal(payload)
				m.Expect(err).To(m.BeNil())
				req, err := http.NewRequest("POST", "/", bytes.NewReader(body))

				router.POST("/", controller.SignUpUser)
				router.ServeHTTP(rr, req)

				m.Expect(rr.Code).To(m.Equal(http.StatusBadRequest))

				response := &generated.BadRequest{}
				err = json.Unmarshal(rr.Body.Bytes(), response)
				m.Expect(err).To(m.BeNil())
				m.Expect(mock.ExpectationsWereMet()).To(m.BeNil())

				m.Expect(response.Detail).To(m.Equal("passwords don't match"))
			})
		})
		g.Context("with existing email", func() {
			g.It("fails with conflict", func() {
				payload := &generated.UserRegister{
					Email:           "user@example.com",
					Name:            "Jason Doe",
					Password:        "password",
					PasswordConfirm: "password",
				}

				mock.ExpectBegin()
				const sqlInsert = `INSERT INTO "users" ("name","email","password","role","verified","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`
				mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
					WithArgs(payload.Name, payload.Email, utils.AnyString{}, "user", true, utils.AnyTime{}, utils.AnyTime{}).
					WillReturnError(fmt.Errorf("duplicate key value violates unique"))

				body, err := json.Marshal(payload)
				m.Expect(err).To(m.BeNil())
				req, err := http.NewRequest("POST", "/", bytes.NewReader(body))

				router.POST("/", controller.SignUpUser)
				router.ServeHTTP(rr, req)

				m.Expect(rr.Code).To(m.Equal(http.StatusConflict))

				response := &generated.ErrorResponse{}
				err = json.Unmarshal(rr.Body.Bytes(), response)
				m.Expect(err).To(m.BeNil())
				m.Expect(mock.ExpectationsWereMet()).To(m.BeNil())

				m.Expect(response.Detail).To(m.Equal("user with that email already exists"))
			})
		})

		g.Context("and sql insert fails with unhandled error", func() {
			g.It("fails with conflict", func() {
				payload := &generated.UserRegister{
					Email:           "user@example.com",
					Name:            "Jason Doe",
					Password:        "password",
					PasswordConfirm: "password",
				}

				mock.ExpectBegin()
				const sqlInsert = `INSERT INTO "users" ("name","email","password","role","verified","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`
				mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
					WithArgs(payload.Name, payload.Email, utils.AnyString{}, "user", true, utils.AnyTime{}, utils.AnyTime{}).
					WillReturnError(fmt.Errorf("something bad happened"))
				mock.ExpectRollback()

				body, err := json.Marshal(payload)
				m.Expect(err).To(m.BeNil())
				req, err := http.NewRequest("POST", "/", bytes.NewReader(body))

				router.POST("/", controller.SignUpUser)
				router.ServeHTTP(rr, req)

				m.Expect(rr.Code).To(m.Equal(http.StatusInternalServerError))

				response := &generated.ErrorResponse{}
				err = json.Unmarshal(rr.Body.Bytes(), response)
				m.Expect(err).To(m.BeNil())
				m.Expect(mock.ExpectationsWereMet()).To(m.BeNil())

				m.Expect(response.Detail).To(m.Equal("failed to create user"))
			})
		})
	})
})
