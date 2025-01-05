package models_test

import (
    "fmt"
    "strings"
    "testing"

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
