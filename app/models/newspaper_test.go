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

type NewspaperTestSuite struct {
    tester.DBSQLiteSuite
    originalDB *gorm.DB
}

func TestNewspaperTestSuite(t *testing.T) {
    suite.Run(t, new(NewspaperTestSuite)) // 対象の構造体にメソッドとして定義されたテストケースを実行できる
}

func (suite *NewspaperTestSuite) SetupSuite() {
    suite.DBSQLiteSuite.SetupSuite()
    suite.originalDB = models.DB // テスト前のデータベースの状態を保存
}

func (suite *NewspaperTestSuite) MockDB() sqlmock.Sqlmock {
    // mockにはsqlmock.Sqlmock（クエリの期待値設定用）が、mockGormDBにはgormのデータベースインスタンスが設定される
    mock, mockGormDB := tester.MockDB()
    models.DB = mockGormDB
    return mock
}

func (suite *NewspaperTestSuite) AfterTest(suiteName, testName string) {
    models.DB = suite.originalDB // テスト前の状態に戻す
}

func (suite *NewspaperTestSuite) TestNewspaper() {
	createdNewspaper, err := models.CreateNewspaper("Test", "sports")
	suite.Assert().Nil(err)
	suite.Assert().Equal("Test", createdNewspaper.Title)
	suite.Assert().Equal("sports", createdNewspaper.ColumnName)

	getNewspaper, err := models.GetNewspaper(createdNewspaper.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("Test", getNewspaper.Title)
	suite.Assert().Equal("sports", getNewspaper.ColumnName)

	getNewspaper.Title = "updated"
	err = getNewspaper.Save()
	suite.Assert().Nil(err)
	updatedNewspaper, err := models.GetNewspaper(createdNewspaper.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("updated", updatedNewspaper.Title)
	suite.Assert().Equal("sports", updatedNewspaper.ColumnName)

	err = updatedNewspaper.Delete()
	suite.Assert().Nil(err)
	deletedNewspaper, err := models.GetNewspaper(updatedNewspaper.ID)
	suite.Assert().Nil(deletedNewspaper)
	suite.Assert().True(strings.Contains("record not found", err.Error()))
}

func (suite *NewspaperTestSuite) TestNewspaperMarshal() {
	newspaper := models.Newspaper{
			Title:      "Test",
			ColumnName: "sports",
	}
	newspaperJSON, err := newspaper.MarshalJSON()
	suite.Assert().Nil(err)
	suite.Assert().JSONEq(fmt.Sprintf(`{
		"columnName":"sports",
		"id":0,
		"title":"Test"
	}`), string(newspaperJSON))
}

func (suite *NewspaperTestSuite) TestNewspaperCreateFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin() // トランザクションの開始を期待
	mockDB.ExpectExec("INSERT INTO `newspapers`").WithArgs("Test", "sports").WillReturnError(errors.New("create error"))
	// トランザクションのロールバックやコミット操作を期待
	mockDB.ExpectRollback()
	mockDB.ExpectCommit()

	newspaper, err := models.CreateNewspaper("Test", "sports")
		suite.Assert().Nil(newspaper)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("create error", err.Error())
}

func (suite *NewspaperTestSuite) TestNewspaperGetFailure() {
	mockDB := suite.MockDB()
	// SQLクエリの期待値を設定。このクエリが実行されると、エラー"get error"が返される
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `newspapers` WHERE `newspapers`.`id` = ? ORDER BY `newspapers`.`id` LIMIT ?")).WithArgs(1, 1).WillReturnError(errors.New("get error"))

	newspaper, err := models.GetNewspaper(1)
	suite.Assert().Nil(newspaper)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}

func (suite *NewspaperTestSuite) TestNewspaperSaveFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin() // トランザクションの開始を期待
	mockDB.ExpectExec(regexp.QuoteMeta("UPDATE `newspapers` SET `title`=?,`column_name`=? WHERE `id` = ?")).WithArgs("updated", "sports", 1).WillReturnError(errors.New("update error"))
	// トランザクションのロールバックやコミット操作を期待
	mockDB.ExpectRollback()

	newspaper := models.Newspaper{
		ID:         1,
		Title:      "Test",
		ColumnName: "sports",
	}
	newspaper.Title = "updated"
	err := newspaper.Save()
	suite.Assert().NotNil(err)
	suite.Assert().Equal("update error", err.Error())
}

func (suite *NewspaperTestSuite) TestNewspaperDeleteFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin() // トランザクションの開始を期待
	mockDB.ExpectExec("DELETE FROM `newspapers` WHERE id = ?").WithArgs(0).WillReturnError(errors.New("delete error"))
	// トランザクションのロールバックやコミット操作を期待
	mockDB.ExpectRollback()
	mockDB.ExpectCommit()

	newspaper := models.Newspaper{
			Title:       "Test",
			ColumnName:  "sports",
	}

	err := newspaper.Delete()
	suite.Assert().NotNil(err)
	suite.Assert().Equal("delete error", err.Error())
}
