package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
		requestBody.NewspaperID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdArticle)
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

fun (a *ArticleHandler) UpdateArticleById(c *gin.Context, ID int) {
	var requestBody api.UpdateArticleByIdJSONRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
		return
	}

	article, err := models.GetArticle(ID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	if requestBody.Body != nil {
		article.Body = *requestBody.Body
	}
	if requestBody.Year != nil {
		article.Year = *requestBody.Year
	}
	if requestBody.Month != nil {
		article.Month = *requestBody.Month
	}
	if requestBody.Day != nil {
		article.Day = *requestBody.Day
	}
	if requestBody.NewspaperID != nil {
		article.NewspaperID = *requestBody.NewspaperID
	}

	if err := article.Save(); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}

fun c (a *ArticleHandler) DeleteArticleById(c *gin.Context, ID int) {
	article := models.Article{ID: ID}

	if err := article.Delete(); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil) // 204 No Content
}
