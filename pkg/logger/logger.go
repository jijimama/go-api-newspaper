package logger

import (
	"os"              // 環境変数の取得やファイル操作のため
	"go.uber.org/zap" // Zapロギングライブラリ
)

// グローバル変数: アプリケーション全体で使用するロガーを定義
var (
	ZapLogger        *zap.Logger       // 構造化ロガー (型安全性の高いLogger)
	zapSugaredLogger *zap.SugaredLogger // 柔軟で使いやすいロガー (SugaredLogger)
)

// 初期化関数: パッケージが読み込まれるときに実行される
func init() {
	// プロダクション環境用のZap設定を生成
	cfg := zap.NewProductionConfig()

	// 環境変数からログファイル名を取得
	logFile := os.Getenv("APP_LOG_FILE")
	if logFile != "" {
		// ログ出力先を標準エラー出力と指定されたファイルに設定
		cfg.OutputPaths = []string{"stderr", logFile}
	}

	// 構造化ロガーを構築 (必須エラー処理でプログラム停止)
	ZapLogger = zap.Must(cfg.Build())

	// 開発環境では、デフォルトの開発用ロガーに切り替え
	if os.Getenv("APP_ENV") == "development" {
		ZapLogger = zap.Must(zap.NewDevelopment()) // 開発用のロガーはより詳細な情報を出力
	}

	// SugaredLoggerを作成して、使いやすさを提供
	zapSugaredLogger = ZapLogger.Sugar()
}

// ログのフラッシュ (重要: ログをファイルや出力先に書き込む)
func Sync() {
	err := zapSugaredLogger.Sync() // メモリに蓄積されたログをすべて出力
	if err != nil {
		// ログ出力に失敗した場合のエラーハンドリング
		zap.Error(err)
	}
}

// Infoレベルのログ出力: 一般的な情報を記録
func Info(msg string, keysAndValues ...interface{}) {
	// メッセージとキー値ペアで構造化ログを出力
	zapSugaredLogger.Infow(msg, keysAndValues...)
}

// Debugレベルのログ出力: デバッグ用の詳細な情報を記録
func Debug(msg string, keysAndValues ...interface{}) {
	// メッセージとキー値ペアでデバッグログを出力
	zapSugaredLogger.Debugw(msg, keysAndValues...)
}

// Warnレベルのログ出力: 警告メッセージを記録
func Warn(msg string, keysAndValues ...interface{}) {
	// メッセージとキー値ペアで警告ログを出力
	zapSugaredLogger.Warnw(msg, keysAndValues...)
}

// Errorレベルのログ出力: エラーメッセージを記録
func Error(msg string, keysAndValues ...interface{}) {
	// メッセージとキー値ペアでエラーログを出力
	zapSugaredLogger.Errorw(msg, keysAndValues...)
}

// Fatalレベルのログ出力: 致命的なエラーを記録し、プログラムを終了
func Fatal(msg string, keysAndValues ...interface{}) {
	// メッセージとキー値ペアで致命的なログを出力
	zapSugaredLogger.Fatalw(msg, keysAndValues...)
}

// Panicレベルのログ出力: パニック状態のエラーを記録し、プログラムを停止
func Panic(msg string, keysAndValues ...interface{}) {
	// メッセージとキー値ペアでパニックログを出力
	zapSugaredLogger.Panicw(msg, keysAndValues...)
}

