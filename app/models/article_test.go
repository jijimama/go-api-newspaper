package models_test

import (
    "errors"
    "fmt"
    "regexp"
    "strings"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/suite"
    "gorm.io/gorm"

    "go-api-newspaper/app/models"
    "go-api-newspaper/pkg/tester"
)

type ArticleTestSuite struct {
	tester.DBSQLiteSuite
	originalDB *gorm.DB
}

func TestArticleTestSuite(t *testing.T) {
	suite.Run(t, new(ArticleTestSuite))
}

func (suite *ArticleTestSuite) SetupSuite() {
    suite.DBSQLiteSuite.SetupSuite()
    suite.originalDB = models.DB
}

func (suite *ArticleTestSuite) MockDB() sqlmock.Sqlmock {
    mock, mockGormDB := tester.MockDB()
    models.DB = mockGormDB
    return mock
}

func (suite *ArticleTestSuite) AfterTest(suiteName, testName string) {
    models.DB = suite.originalDB
}

func (suite *ArticleTestSuite) TestArticle() {
    createdArticle, err := models.CreateArticle("test body", 2023, 10, 2, 1)
    suite.Assert().Nil(err)
    suite.Assert().Equal("test body", createdArticle.Body)
    suite.Assert().Equal(2023, createdArticle.Year)
    suite.Assert().Equal(10, createdArticle.Month)
    suite.Assert().Equal(2, createdArticle.Day)
    suite.Assert().Equal(1, createdArticle.NewspaperID)
}
