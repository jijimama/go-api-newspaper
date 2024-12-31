package models

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go-api-newspaper/configs"
)

const (
	InstanceSqlLite int = iota // SQLiteを選択する定数
	InstanceMySQL              // MySQLを選択する定数
)

var (
	DB                            *gorm.DB // グローバルで利用可能なデータベースインスタンス
	errInvalidSQLDatabaseInstance = errors.New("invalid sql db instance") // 不正なデータベースインスタンスを扱うエラー
)

// 初期化するモデル（テーブル）をリストで返す関数
func GetModels() []interface{} {
	return []interface{}{&Newspaper{}}
}

// データベースのインスタンスを生成するファクトリ関数
func NewDatabaseSQLFactory(instance int) (db *gorm.DB, err error) {
	switch instance {
	case InstanceMySQL:
		// MySQL用の接続情報（DSN: Data Source Name）を構築
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
			configs.Config.DBUser,
			configs.Config.DBPassword,
			configs.Config.DBHost,
			configs.Config.DBPort,
			configs.Config.DBName)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case InstanceSqlLite:
		db, err = gorm.Open(sqlite.Open(configs.Config.DBName), &gorm.Config{})
	default:
		return nil, errInvalidSQLDatabaseInstance
	}
	return db, err
}

// データベースをセットする関数（実際にグローバル変数DBにインスタンスを格納）
func SetDatabase(instance int) (err error) {
	db, err := NewDatabaseSQLFactory(instance)
	if err != nil {
		return err
	}
	DB = db // グローバル変数にセット
	DB = db
	return err
}
