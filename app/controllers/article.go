package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"go-api-newspaper/api"
	"go-api-newspaper/app/models"
	"go-api-newspaper/pkg/logger"
)

type ArticleHandler struct{}

func (a *ArticleHandler) CreateArticle(c *gin.Context) {
	var requestBody api.CreateArticleJSONRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
		return
	}

	createdArticle, err := models.CreateArticle(
		requestBody.Body,
		requestBody.Year,
		requestBody.Month,
		requestBody.Day,
		*requestBody.NewspaperID,
	)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdArticle) // 201 レスポンスに書き込み
}

func (a *ArticleHandler) GetArticleById(c *gin.Context, ID int) {
	article, err := models.GetArticle(ID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}
