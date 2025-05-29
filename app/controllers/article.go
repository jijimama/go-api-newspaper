package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"go-api-newspaper/api"
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
}
