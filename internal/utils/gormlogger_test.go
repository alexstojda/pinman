package utils_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	g "github.com/onsi/ginkgo/v2"
	m "github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/logger"
	"pinman/internal/utils"
	"time"
)

var _ = g.Describe("GormLogger", func() {
	var gormLogger utils.GormLogger
	var logr zerolog.Logger
	var ctx context.Context
	var b *bytes.Buffer

	g.BeforeEach(func() {
		b = &bytes.Buffer{}
		gormLogger = utils.GormLogger{}
		logr = log.With().Logger().Output(b)
		ctx = logr.WithContext(context.TODO())
	})

	g.When("LogMode is called", func() {
		g.It("returns the logger unchanged", func() {
			logModeRes := gormLogger.LogMode(logger.Info)
			m.Expect(gormLogger).To(m.BeEquivalentTo(logModeRes))
		})
	})

	g.When("Error is called", func() {
		g.It("should run without error", func() {
			gormLogger.Error(ctx, "an %s", "error")
			m.Expect(b.String()).To(m.MatchRegexp("^{\"level\":\"error\".*"))
		})
	})

	g.When("Warn is called", func() {
		g.It("should run without error", func() {
			gormLogger.Warn(ctx, "an %s", "error")
			m.Expect(b.String()).To(m.MatchRegexp("^{\"level\":\"warn\".*"))
		})
	})

	g.When("Info is called", func() {
		g.It("should run without error", func() {
			gormLogger.Info(ctx, "an %s", "error")
			m.Expect(b.String()).To(m.MatchRegexp("^{\"level\":\"info\".*"))
		})
	})

	g.When("Trace is called", func() {
		g.It("should run", func() {
			t := time.Now()
			gormLogger.Trace(ctx, t, func() (string, int64) {
				return "SELECT * FROM users", 69
			}, errors.New("foobar"))
			m.Expect(b.String()).To(
				m.MatchRegexp(fmt.Sprintf(
					`{"level":"debug",".*":\d{1,6}(\.\d{1,25})?,"sql":"SELECT \* FROM users","rows":69,"time":"%s"}`,
					t.Format(time.RFC3339),
				)),
			)
		})
		g.Context("with non-default duration field", func() {
			var unit time.Duration
			g.It("shows elapsed in nanoseconds", func() {
				unit = time.Nanosecond
			})
			g.It("shows elapsed in microseconds", func() {
				unit = time.Microsecond
			})
			g.It("shows elapsed in millisecond", func() {
				unit = time.Millisecond
			})
			g.It("shows elapsed in seconds", func() {
				unit = time.Second
			})
			g.It("shows elapsed in minutes", func() {
				unit = time.Minute
			})
			g.It("shows elapsed in hours", func() {
				unit = time.Hour
			})
			g.It("shows elapsed in non-standard duration", func() {
				unit = 12 * time.Minute
			})
			g.AfterEach(func() {
				zerolog.DurationFieldUnit = unit
				t := time.Now()
				gormLogger.Trace(ctx, t, func() (string, int64) {
					return "SELECT * FROM users", 69
				}, nil)
				m.Expect(b.String()).To(
					m.MatchRegexp(fmt.Sprintf(
						`{"level":"trace",".*":\d{1,6}(\.\d{1,30})?,"sql":"SELECT \* FROM users","rows":69,"time":"%s"}`,
						t.Format(time.RFC3339),
					)),
				)
				// Reset back to default
				zerolog.DurationFieldUnit = time.Millisecond
			})
		})
	})
})
