package tester

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"go-api-newspaper/app/models"
	"go-api-newspaper/configs"
)

// testcontainersを利用するには、Dockerが起動している必要
// CheckPortは指定されたホストとポートに接続可能かを確認する関数
func CheckPort(host string, port int) bool {
	// 指定されたホストとポートにTCP接続を試みる
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if conn != nil {
		// 接続が成功した場合は閉じてfalseを返す（ポートが使用中）
		conn.Close()
		return false
	}
	if err != nil {
		// 接続が失敗した場合はtrueを返す（ポートが空いている）
		return true
	}
	return false
}

// WaitForPortは、指定したホストとポートが空くのを待つ関数
func WaitForPort(host string, port int, timeout time.Duration) bool {
	// 現在の時間にタイムアウトを加えた時刻を計算
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		// ポートが空いているかを1秒ごとにチェック
		if CheckPort(host, port) {
			return true
		}
		time.Sleep(1 * time.Second)
	}
	// タイムアウトが経過しても空かなかった場合はfalseを返す
	return false
}

// Mysqlに接続するための構造体
type DBMySQLSuite struct {
	suite.Suite                             // testifyのSuite機能を埋め込む
	mySQLContainer testcontainers.Container // MySQLコンテナのインスタンス
	ctx            context.Context          // コンテナ操作用のコンテキスト
}

// SetupTestContainersは、MySQLコンテナをセットアップする関数
func (suite *DBMySQLSuite) SetupTestContainers() (err error) {
	// Dockerがポートを使用できるまで待機
	WaitForPort(configs.Config.DBHost, configs.Config.DBPort, 10*time.Second)
	suite.ctx = context.Background() // コンテキストを初期化
	req := testcontainers.ContainerRequest{ // コンテナに対する設定のリクエストを作成
		Image: "mysql:8",
		Env: map[string]string{
			"MYSQL_DATABASE":             configs.Config.DBName,
			"MYSQL_USER":                 configs.Config.DBUser,
			"MYSQL_PASSWORD":             configs.Config.DBPassword,
			"MYSQL_ALLOW_EMPTY_PASSWORD": "yes",
		},
		ExposedPorts: []string{fmt.Sprintf("%d:3306/tcp", configs.Config.DBPort)}, // ポートマッピング
		WaitingFor:   wait.ForLog("port: 3306  MySQL Community Server"),           // MySQLの準備完了をログで確認
	}
	// リクエストをもとにコンテナを作成
	suite.mySQLContainer, err = testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

// SetupSuiteはテストスイート全体の初期設定を行う関数
func (suite *DBMySQLSuite) SetupSuite() {
	// MySQLコンテナのセットアップを実行
	err := suite.SetupTestContainers()
	suite.Assert().Nil(err)

	// モデルにMySQLデータベースを設定
	err = models.SetDatabase(models.InstanceMySQL)
	suite.Assert().Nil(err)

	for _, model := range models.GetModels() {
		err := models.DB.AutoMigrate(model)
		suite.Assert().Nil(err)
	}
}

// TearDownSuiteはテストスイート全体のクリーンアップを行う関数
func (suite *DBMySQLSuite) TearDownSuite() {
	if suite.mySQLContainer == nil {
		return
	}
	// MySQLコンテナを終了（停止と削除）する
	err := suite.mySQLContainer.Terminate(suite.ctx)
	suite.Assert().Nil(err)
}
