package tester

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-api-newspaper/pkg/logger"
)

func MockDB() (mock sqlmock.Sqlmock, mockGormDB *gorm.DB) {
	mockDB, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		logger.Fatal(err.Error())
	}

	mockGormDB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       "mock_db",
		DriverName:                "mysql",
		Conn:                      mockDB, // MySQLドライバが実際のデータベース接続の代わりにモック接続を使用
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		logger.Fatal(err.Error())
	}
	return mock, mockGormDB
}

type mockClock struct {
	t time.Time
}

func NewMockClock(t time.Time) mockClock {
	return mockClock{t}
}

func (m mockClock) Now() time.Time {
	return m.t
}
