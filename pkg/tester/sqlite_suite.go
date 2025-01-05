package tester

import (
	"os"

	"github.com/stretchr/testify/suite"

	"go-api-newspaper/app/models"
	"go-api-newspaper/configs"
)

type DBSQLiteSuite struct {
	suite.Suite
}

// テスト前に自動で実行されるメソッド
func (suite *DBSQLiteSuite) SetupSuite() {
	configs.Config.DBName = "unittest.sqlite"
	err := models.SetDatabase(models.InstanceSqlLite) // 初期化（unittest.sqliteというデータベースが保存）
	suite.Assert().Nil(err)

	for _, model := range models.GetModels() {
		err := models.DB.AutoMigrate(model) // モデルの構造に対応したテーブルを作成
		suite.Assert().Nil(err)             // エラーがないこと確認
	}
}

// テスト後に実行されるメソッド
func (suite *DBSQLiteSuite) TearDownSuite() {
	err := os.Remove(configs.Config.DBName) // データベースファイルを削除
	suite.Assert().Nil(err)
}
