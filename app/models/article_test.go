package models_test

import (
	"gorm.io/gorm"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"go-api-newspaper/app/models"
	"go-api-newspaper/pkg/tester"
)

type ArticleTestSuite struct {
	tester.DBSQLiteSuite
	originalDB *gorm.DB
}

func TestArticleTestSuite(t *testing.T) {
	suite.Run(t, new(ArticleTestSuite)) // 対象の構造体にメソッドとして定義されたテストケースを実行できる
}

func (suite *ArticleTestSuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.originalDB = models.DB // テスト前のデータベースの状態を保存
}

func (suite *ArticleTestSuite) MockDB() sqlmock.Sqlmock {
	// mockにはsqlmock.Sqlmock（クエリの期待値設定用）が、mockGormDBにはgormのデータベースインスタンスが設定される
	mock, mockGormDB := tester.MockDB()
	models.DB = mockGormDB
	return mock
}

func (suite *ArticleTestSuite) AfterTest(suiteName, testName string) {
	models.DB = suite.originalDB // テスト前の状態に戻す
}

func (suite *ArticleTestSuite) TestArticle() {
	createdArticle, err := models.CreateArticle("Test", 2023, 10, 1, 1)
	suite.Assert().Nil(err)
	suite.Assert().Equal("Test", createdArticle.Body)
}
