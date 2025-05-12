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
	err := models.DB.AutoMigrate(&models.Article{}, &models.Newspaper{})
	suite.Assert().Nil(err, "マイグレーションに失敗しました")
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
	createdNewspaper, err := models.CreateNewspaper("Test Newspaper", "Test Column")
	suite.Assert().Nil(err)

	createdArticle, err := models.CreateArticle("Test", 2023, 10, 1, createdNewspaper.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("Test", createdArticle.Body)
	suite.Assert().Equal(2023, createdArticle.Year)
	suite.Assert().Equal(10, createdArticle.Month)
	suite.Assert().Equal(1, createdArticle.Day)
	suite.Assert().Equal(1, createdArticle.NewspaperID)

	getArticle, err := models.GetArticle(createdArticle.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("Test", getArticle.Body)
	suite.Assert().Equal(2023, getArticle.Year)
	suite.Assert().Equal(10, getArticle.Month)
	suite.Assert().Equal(1, getArticle.Day)
	suite.Assert().Equal(1, getArticle.NewspaperID)

	getArticle.Body = "updated"
	err = getArticle.Save()
	suite.Assert().Nil(err)
	updatedArticle, err := models.GetArticle(createdArticle.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("updated", updatedArticle.Body)
	suite.Assert().Equal(2023, updatedArticle.Year)
	suite.Assert().Equal(10, updatedArticle.Month)
	suite.Assert().Equal(1, updatedArticle.Day)
	suite.Assert().Equal(1, updatedArticle.NewspaperID)
}
