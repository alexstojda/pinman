package utils

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http/httptest"
	"time"

	m "github.com/onsi/gomega"
)

func NewGormMock() (db *gorm.DB, mock sqlmock.Sqlmock) {
	sqlDb, mock, err := sqlmock.New()
	m.Expect(err).To(m.BeNil())
	m.Expect(mock).ToNot(m.BeNil())
	m.Expect(sqlDb).ToNot(m.BeNil())

	db, err = ConnectDB(&Config{}, postgres.New(postgres.Config{
		Conn: sqlDb,
	}))
	m.Expect(err).To(m.BeNil())

	return db, mock
}

func NewGinTestCtx() (*gin.Context, *httptest.ResponseRecorder, *gin.Engine) {
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, router := gin.CreateTestContext(recorder)

	return ctx, recorder, router
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type AnyString struct{}

func (a AnyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}
