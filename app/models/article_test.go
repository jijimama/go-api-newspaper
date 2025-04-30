package models_test

import (
	"gorm.io/gorm"

	"go-api-newspaper/pkg/tester"
)

type ArticleTestSuite struct {
	tester.DBSQLiteSuite
	originalDB *gorm.DB
}
