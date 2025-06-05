package models_test

import (
	"errors"
	"gorm.io/gorm"
	"regexp"
	"strings"
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

	err = updatedArticle.Delete()
	suite.Assert().Nil(err)
	deletedArticle, err := models.GetArticle(updatedArticle.ID)
	suite.Assert().Nil(deletedArticle)
	suite.Assert().True(err != nil)
	suite.Assert().True(strings.Contains("record not found", err.Error()))
}

func (suite *ArticleTestSuite) TestArticleMarshal() {
	article := models.Article{
		ID:          1,
		Body:        "Test",
		Year:        2023,
		Month:       10,
		Day:         1,
		NewspaperID: 1,
		Newspaper: &models.Newspaper{
			ID:         1,
			Title:      "Test Newspaper",
			ColumnName: "Test Column",
		},
	}
	newspaperJSON, err := article.MarshalJSON()
	suite.Assert().Nil(err)
	suite.Assert().Equal(`{"body":"Test","day":1,"id":1,"month":10,"newspaperID":1,"year":2023}`, string(newspaperJSON))
}

func (suite *ArticleTestSuite) TestArticleCreateFailure() {
	mockDB := suite.MockDB()

	newspaper := models.Newspaper{
		ID:         1,
		Title:      "Test",
		ColumnName: "sports",
	}

	// 新聞の取得をモック
	mockDB.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `newspapers` WHERE `newspapers`.`id` = ? ORDER BY `newspapers`.`id` LIMIT ?",
	)).WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "column_name"}).
		AddRow(newspaper.ID, newspaper.Title, newspaper.ColumnName),
	)

	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `newspapers` (`title`,`column_name`,`id`) VALUES (?,?,?) ON DUPLICATE KEY UPDATE `id`=`id`",
	)).WithArgs(newspaper.Title, newspaper.ColumnName, newspaper.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockDB.ExpectExec(regexp.QuoteMeta("INSERT INTO `articles`")).
		WithArgs("Test", 2023, 10, 1, newspaper.ID).
		WillReturnError(errors.New("create error"))

	mockDB.ExpectRollback()

	article, err := models.CreateArticle("Test", 2023, 10, 1, newspaper.ID)

	suite.Assert().Nil(article)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("create error", err.Error())
}

func (suite *ArticleTestSuite) TestArticleGetFailure() {
	// MockDBを使用して、Getメソッドがエラーを返すように設定
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `articles` WHERE id = ? ORDER BY `articles`.`id` LIMIT ?")).WithArgs(1, 1).WillReturnError(errors.New("get error"))

	article, err := models.GetArticle(1)
	suite.Assert().Nil(article)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}

func (suite *ArticleTestSuite) TestArticleSaveFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta(
		"UPDATE `articles` SET `body`=?,`year`=?,`month`=?,`day`=?,`newspaper_id`=? WHERE `id` = ?",
	)).WithArgs("updated", 2023, 10, 1, 1, 1).
		WillReturnError(errors.New("update error"))

	mockDB.ExpectRollback()

	article := models.Article{
		ID:          1,
		Body:        "Test",
		Year:        2023,
		Month:       10,
		Day:         1,
		NewspaperID: 1,
	}

	article.Body = "updated"

	err := article.Save()
	suite.Assert().NotNil(err)
	suite.Assert().Equal("update error", err.Error())
}

func (suite *ArticleTestSuite) TestArticleDeleteFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta(
		"DELETE FROM `articles` WHERE id = ?",
	)).WithArgs(0).
		WillReturnError(errors.New("delete error"))

	mockDB.ExpectRollback()

	article := models.Article{
		ID:          0,
		Body:        "Test",
		Year:        2023,
		Month:       10,
		Day:         1,
		NewspaperID: 1,
	}

	err := article.Delete()
	suite.Assert().NotNil(err)
	suite.Assert().Equal("delete error", err.Error())
}
