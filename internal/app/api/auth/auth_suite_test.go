package auth_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"net/url"
	"pinman/internal/app/api/auth"
	"pinman/internal/app/generated"
	"pinman/internal/models"
	"pinman/internal/utils"
	"regexp"
	"strings"
	"testing"
	"time"

	g "github.com/onsi/ginkgo/v2"
	m "github.com/onsi/gomega"
)

func TestAuth(t *testing.T) {
	m.RegisterFailHandler(g.Fail)
	g.RunSpecs(t, "Auth Suite")
}

var _ = g.Describe("JWTMiddleware", func() {
	config := &utils.Config{
		TokenSecretKey: "asecretkey",
	}

	g.When("CreateJWTMiddleware is called", func() {
		var db *gorm.DB

		g.BeforeEach(func() {
			db, _ = utils.NewGormMock()
		})

		g.Context("with valid parameters", func() {
			mw, err := auth.CreateJWTMiddleware(
				config,
				db,
			)

			g.It("should succeed", func() {
				m.Expect(err).To(m.BeNil())
				m.Expect(mw).ToNot(m.BeNil())
			})
		})
	})

	g.When("LoginHandler receives a request", func() {
		var db *gorm.DB
		var mock sqlmock.Sqlmock
		var rr *httptest.ResponseRecorder
		var router *gin.Engine
		var mw *jwt.GinJWTMiddleware

		var payload = &generated.UserLogin{
			Username: "email@example.com",
			Password: "password",
		}
		var user *models.User
		var columns = []string{
			"id", "name", "email", "created_at", "updated_at", "password", "role", "verified",
		}
		var rows *sqlmock.Rows

		g.BeforeEach(func() {
			db, mock = utils.NewGormMock()
			_, rr, router = utils.NewGinTestCtx()
			var err error
			mw, err = auth.CreateJWTMiddleware(config, db)
			if err != nil {
				g.Fail("could not create middleware")
			}

			uid, err := uuid.NewUUID()
			m.Expect(err).To(m.BeNil())
			hashedPass, err := utils.HashPassword("password")
			m.Expect(err).To(m.BeNil())
			user = &models.User{
				ID:        uid,
				Name:      "John Doe",
				Email:     "email@example.com",
				Password:  hashedPass,
				Role:      "user",
				Verified:  true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			rows = sqlmock.NewRows(columns).AddRow(
				user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt, user.Password, user.Role, user.Verified,
			)
		})

		g.It("should be successful for json request with valid credentials", func() {
			mock.ExpectQuery(
				regexp.QuoteMeta(
					`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT 1`,
				),
			).WithArgs(strings.ToLower(payload.Username)).
				WillReturnRows(rows)

			body, err := json.Marshal(payload)
			m.Expect(err).To(m.BeNil())
			req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
			m.Expect(err).To(m.BeNil())

			req.Header.Set("Content-Type", "application/json")
			router.POST("/", mw.LoginHandler)
			router.ServeHTTP(rr, req)

			m.Expect(rr.Code).To(m.Equal(http.StatusOK))
			m.Expect(mock.ExpectationsWereMet()).To(m.BeNil())

			response := generated.TokenResponse{}
			m.Expect(json.Unmarshal(rr.Body.Bytes(), &response)).To(m.BeNil())
			m.Expect(response.AccessToken).ToNot(m.BeEmpty())
			m.Expect(response.Expire).ToNot(m.BeEmpty())
		})

		g.It("should be successful for urlencoded request with valid credentials", func() {
			mock.ExpectQuery(
				regexp.QuoteMeta(
					`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT 1`,
				),
			).WithArgs(strings.ToLower(payload.Username)).
				WillReturnRows(rows)

			body := url.Values{}
			body.Set("username", payload.Username)
			body.Set("password", payload.Password)

			req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(body.Encode())))
			m.Expect(err).To(m.BeNil())

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			router.POST("/", mw.LoginHandler)
			router.ServeHTTP(rr, req)

			m.Expect(rr.Code).To(m.Equal(http.StatusOK))
			m.Expect(mock.ExpectationsWereMet()).To(m.BeNil())

			response := generated.TokenResponse{}
			m.Expect(json.Unmarshal(rr.Body.Bytes(), &response)).To(m.BeNil())
			m.Expect(response.AccessToken).ToNot(m.BeEmpty())
			m.Expect(response.Expire).ToNot(m.BeEmpty())
		})

		g.It("should return as unauthorized for invalid password", func() {
			payload.Password = "wrong password"
			mock.ExpectQuery(
				regexp.QuoteMeta(
					`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT 1`,
				),
			).WithArgs(strings.ToLower(payload.Username)).
				WillReturnRows(rows)

			body, err := json.Marshal(payload)
			m.Expect(err).To(m.BeNil())
			req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
			m.Expect(err).To(m.BeNil())

			req.Header.Set("Content-Type", "application/json")
			router.POST("/", mw.LoginHandler)
			router.ServeHTTP(rr, req)

			m.Expect(rr.Code).To(m.Equal(http.StatusUnauthorized))
			m.Expect(mock.ExpectationsWereMet()).To(m.BeNil())

			response := generated.ErrorResponse{}
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			m.Expect(err).To(m.BeNil())
			m.Expect(response.Status).To(m.BeEquivalentTo(http.StatusUnauthorized))
		})

		g.It("should return as unauthorized for unknown email", func() {
			payload.Username = "unknown@example.com"
			mock.ExpectQuery(
				regexp.QuoteMeta(
					`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT 1`,
				),
			).WithArgs(strings.ToLower(payload.Username)).
				WillReturnRows(sqlmock.NewRows(columns))

			body, err := json.Marshal(payload)
			m.Expect(err).To(m.BeNil())
			req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
			m.Expect(err).To(m.BeNil())
			req.Header.Set("Content-Type", "application/json")

			router.POST("/", mw.LoginHandler)
			router.ServeHTTP(rr, req)

			m.Expect(rr.Code).To(m.Equal(http.StatusUnauthorized))
			m.Expect(mock.ExpectationsWereMet()).To(m.BeNil())

			response := generated.ErrorResponse{}
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			m.Expect(err).To(m.BeNil())
			m.Expect(response.Status).To(m.BeEquivalentTo(http.StatusUnauthorized))
		})

		g.It("should return as unauthorized for invalid request body", func() {
			payload := map[string]string{
				"Foo": "bar",
			}

			body, err := json.Marshal(payload)
			m.Expect(err).To(m.BeNil())
			req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
			m.Expect(err).To(m.BeNil())

			req.Header.Set("Content-Type", "application/json")
			router.POST("/", mw.LoginHandler)
			router.ServeHTTP(rr, req)

			m.Expect(rr.Code).To(m.Equal(http.StatusUnauthorized))
			m.Expect(mock.ExpectationsWereMet()).To(m.BeNil())

			response := generated.Unauthorized{}
			m.Expect(json.Unmarshal(rr.Body.Bytes(), &response)).To(m.BeNil())
		})

		g.It("should return as unauthorized for invalid content-type", func() {
			req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte{}))
			m.Expect(err).To(m.BeNil())

			req.Header.Set("Content-Type", "application/xml")
			router.POST("/", mw.LoginHandler)
			router.ServeHTTP(rr, req)

			m.Expect(rr.Code).To(m.Equal(http.StatusUnauthorized))
			m.Expect(mock.ExpectationsWereMet()).To(m.BeNil())

			response := generated.Unauthorized{}
			m.Expect(json.Unmarshal(rr.Body.Bytes(), &response)).To(m.BeNil())
		})
	})

	g.When("RefreshHandler receives a request", func() {
		var db *gorm.DB
		var rr *httptest.ResponseRecorder
		var router *gin.Engine
		var mw *jwt.GinJWTMiddleware
		var user *models.User
		var token string

		g.BeforeEach(func() {
			db, _ = utils.NewGormMock()
			_, rr, router = utils.NewGinTestCtx()
			var err error
			mw, err = auth.CreateJWTMiddleware(config, db)
			m.Expect(err).To(m.BeNil())

			uid, err := uuid.NewUUID()
			m.Expect(err).To(m.BeNil())
			hashedPass, err := utils.HashPassword("password")
			m.Expect(err).To(m.BeNil())
			user = &models.User{
				ID:        uid,
				Name:      "John Doe",
				Email:     "email@example.com",
				Password:  hashedPass,
				Role:      "user",
				Verified:  true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		})

		g.Context("with valid token", func() {
			g.It("should be successful", func() {
				var err error
				token, _, err = mw.TokenGenerator(user)
				m.Expect(err).To(m.BeNil())

				req, err := http.NewRequest("GET", "/", nil)
				m.Expect(err).To(m.BeNil())
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

				router.GET("/", mw.RefreshHandler)
				router.ServeHTTP(rr, req)

				m.Expect(rr.Code).To(m.Equal(http.StatusOK))

				response := generated.TokenResponse{}
				m.Expect(json.Unmarshal(rr.Body.Bytes(), &response)).To(m.BeNil())
				m.Expect(response.AccessToken).ToNot(m.BeEmpty())
				m.Expect(response.Expire).ToNot(m.BeEmpty())
			})
		})
	})

	g.When("MiddlewareFunc receives a request", func() {
		var db *gorm.DB
		var mock sqlmock.Sqlmock
		//var ctx *gin.Context
		var rr *httptest.ResponseRecorder
		var router *gin.Engine
		var mw *jwt.GinJWTMiddleware
		var user *models.User

		var columns = []string{
			"id", "name", "email", "created_at", "updated_at", "password", "role", "verified",
		}
		var rows *sqlmock.Rows

		g.BeforeEach(func() {
			db, mock = utils.NewGormMock()
			_, rr, router = utils.NewGinTestCtx()
			var err error
			mw, err = auth.CreateJWTMiddleware(config, db)
			if err != nil {
				g.Fail("could not create middleware")
			}

			uid, err := uuid.NewUUID()
			m.Expect(err).To(m.BeNil())
			hashedPass, err := utils.HashPassword("password")
			m.Expect(err).To(m.BeNil())
			user = &models.User{
				ID:        uid,
				Name:      "John Doe",
				Email:     "email@example.com",
				Password:  hashedPass,
				Role:      "user",
				Verified:  true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			rows = sqlmock.NewRows(columns).AddRow(
				user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt, user.Password, user.Role, user.Verified,
			)
		})

		g.Context("that does not require authentication", func() {
			g.It("does nothing", func() {
				req, err := http.NewRequest("GET", "/", bytes.NewReader([]byte{}))
				m.Expect(err).To(m.BeNil())

				router.GET("/", func(context *gin.Context) {
					auth.GetAuthMiddlewareFunc(mw)(context)
					m.Expect(context.IsAborted()).To(m.BeFalse())
				})
				router.ServeHTTP(rr, req)
			})
		})

		g.Context("that requires authentication", func() {
			g.It("succeeds if user has required auth scopes", func() {
				mock.ExpectQuery(
					regexp.QuoteMeta(
						`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT 1`,
					),
				).WithArgs(strings.ToLower(user.ID.String())).
					WillReturnRows(rows)

				token, _, err := mw.TokenGenerator(user)

				req, err := http.NewRequest("GET", "/", bytes.NewReader([]byte{}))
				m.Expect(err).To(m.BeNil())
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

				router.GET("/", func(context *gin.Context) {
					context.Set(generated.PinmanAuthScopes, []string{"user"})
					auth.GetAuthMiddlewareFunc(mw)(context)
					m.Expect(context.IsAborted()).To(m.BeFalse())
				})
				router.ServeHTTP(rr, req)
			})
			g.It("fails if user does not have required auth scope", func() {
				mock.ExpectQuery(
					regexp.QuoteMeta(
						`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT 1`,
					),
				).WithArgs(strings.ToLower(user.ID.String())).
					WillReturnRows(rows)

				token, _, err := mw.TokenGenerator(user)

				req, err := http.NewRequest("GET", "/", bytes.NewReader([]byte{}))
				m.Expect(err).To(m.BeNil())
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

				router.GET("/", func(context *gin.Context) {
					context.Set(generated.PinmanAuthScopes, []string{"admin"})
					auth.GetAuthMiddlewareFunc(mw)(context)
					m.Expect(context.IsAborted()).To(m.BeTrue())
				})
				router.ServeHTTP(rr, req)
			})
		})
	})
})
