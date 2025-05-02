package models_test

import (
	"gorm.io/gorm"
	"testing"

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
