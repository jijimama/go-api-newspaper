package configs

import (
	"os"               // 環境変数の取得に使用。
	"strconv"          // 文字列を数値に変換するため。

	"go.uber.org/zap"  //高速で構造化されたロギングライブラリ。

	"go-api-newspaper/pkg/logger"
)

func GetEnvDefault(key, defVal string) string {
	val, err := os.LookupEnv(key) // 指定したキーの環境変数を探し、値と存在の有無を返します。
	if !err { // 環境変数が存在しない場合にデフォルト値を使用。
		return defVal
	}
	return val
}

// ConfigList 構造体は、アプリケーションの設定を保持します。
type ConfigList struct {
	Env                 string
	DBHost              string
	DBPort              int
	DBDriver            string
	DBName              string
	DBUser              string
	DBPassword          string
	APICorsAllowOrigins []string
}

// 環境が開発用かどうかを判定するメソッド
func (c *ConfigList) IsDevelopment() bool {
	return c.Env == "development"
}

var Config ConfigList

// LoadEnv 関数は、環境変数を読み込み、設定を構築します。
func LoadEnv() error { // MYSQL_PORT は数値変換が必要なので、strconv.Atoi() を使用。
	DBPort, err := strconv.Atoi(GetEnvDefault("MYSQL_PORT", "3306"))
	if err != nil {
		return err
	}

	Config = ConfigList{
		Env:                 GetEnvDefault("APP_ENV", "development"),
		DBDriver:            GetEnvDefault("DB_DRIVER", "mysql"),
		DBHost:              GetEnvDefault("DB_HOST", "0.0.0.0"),
		DBPort:              DBPort,
		DBUser:              GetEnvDefault("DB_USER", "app"),
		DBPassword:          GetEnvDefault("DB_PASSWORD", "password"),
		DBName:              GetEnvDefault("DB_NAME", "api_database"),
		APICorsAllowOrigins: []string{"http://0.0.0.0:8001"},
	}
	return nil
}

func init() {
	if err := LoadEnv(); err != nil {
		logger.Error("Failed to load env: ", zap.Error(err))
		panic(err)
	}
}
