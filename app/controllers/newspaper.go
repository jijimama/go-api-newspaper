package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-api-newspaper/api"
	"go-api-newspaper/app/models"
	"go-api-newspaper/pkg/logger"
)
// メソッドを関連付けることで、各エンドポイントの処理を実装。
type NewspaperHandler struct{}

func (a *NewspaperHandler) CreateNewspaper(c *gin.Context) { //*gin.Context リクエストやレスポンスの情報を保持する
	var requestBody api.CreateNewspaperJSONRequestBody         // 自動生成済み
	if err := c.ShouldBindJSON(&requestBody); err != nil { // JSONリクエストボディを構造体にバインド（マッピング）
		logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
		return
	}

	createdNewspaper, err := models.CreateNewspaper(
		requestBody.Title,
		requestBody.ColumnName)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdNewspaper) // 201 レスポンスに書き込み
}

func (a *NewspaperHandler) GetNewspaperById(c *gin.Context, ID int) {
	newspaper, err := models.GetNewspaper(ID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, newspaper)
}

func (a *NewspaperHandler) UpdateNewspaperById(c *gin.Context, ID int) {
	var requestBody api.UpdateNewspaperByIdJSONRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil { // 引数cの内容をrequestBodyに格納
		logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
		return
	}

	newspaper, err := models.GetNewspaper(ID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	if requestBody.Title != nil {
		newspaper.Title = *requestBody.Title
	}
	if requestBody.ColumnName != nil {
		newspaper.ColumnName = *requestBody.ColumnName
	}

	if err := newspaper.Save(); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, newspaper)
}

func (a *NewspaperHandler) DeleteNewspaperById(c *gin.Context, ID int) {
	newspaper := models.Newspaper{ID: ID}

	if err := newspaper.Delete(); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil) // 204
}
