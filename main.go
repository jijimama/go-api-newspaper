package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-api-newspaper/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"

	"go-api-newspaper/api"
	"go-api-newspaper/app/controllers"
	"go-api-newspaper/app/models"
	"go-api-newspaper/configs"
)

func corsMiddleware(allowOrigins []string) gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = allowOrigins
	return cors.New(config)
}

func timeoutMiddleware(duration time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(duration),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(func(c *gin.Context) {
			c.JSON(
				http.StatusRequestTimeout,
				api.ErrorResponse{Message: "timeout"},
			)
			c.Abort()
		}),
	)
}

func main() {
	if err := models.SetDatabase(models.InstanceMySQL); err != nil {
		logger.Fatal(err.Error())
	}

	router := gin.Default() // HTTPリクエストを振り分けるためのルーター

	router.Use(corsMiddleware(configs.Config.APICorsAllowOrigins))

	// OpenAPI仕様を取得（API仕様のバリデーション用）
	swagger, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}

	// 開発環境の場合、Swagger UIを有効化
	if configs.Config.IsDevelopment() {
		swaggerJson, _ := json.Marshal(swagger) // OpenAPI仕様をJSON形式に変換
		var SwaggerInfo = &swag.Spec{
			InfoInstanceName: "swagger",
			SwaggerTemplate:  string(swaggerJson),
		}
		swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo) // Swagger情報を登録
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler)) // Swagger UIエンドポイントを追加
	}

	router.GET("/health", controllers.Health)
	
	apiGroup := router.Group("/api")
	{
		apiGroup.Use(timeoutMiddleware(2 * time.Second))
		v1 := apiGroup.Group("/v1")
		{
			// OpenAPI仕様に基づくリクエストバリデーションをミドルウェアとして追加
			v1.Use(middleware.OapiRequestValidator(swagger)) // 変数swaggerのAPI仕様に基づくバリデーション
			newspaperHandler := &controllers.NewspaperHandler{}
			api.RegisterHandlers(v1, newspaperHandler) // ルーターに登録
		}
	}

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	go func() { // ListenAndServe　サーバーを起動しリクエストをまつ
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)                      //  os.Signalを受け取るためのチャネルを作成
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // シグナルがあればquitチャネルに送信
	<-quit
	log.Println("Shutdown Server ...")
	defer logger.Sync() // ログのバッファをフラッシュする

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // ２秒のタイムアウトを持つコンテキスト
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Server Shutdown: %s", err.Error()))
	}
	<-ctx.Done()
	logger.Info("Shutdown")
}
