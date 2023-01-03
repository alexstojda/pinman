package utils

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	m "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http/httptest"
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
