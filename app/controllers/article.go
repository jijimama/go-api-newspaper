package controllers

import (
	"github.com/gin-gonic/gin"

	"go-api-newspaper/api"
)

type ArticleHandler struct{}

func (a *ArticleHandler) CreateArticle(c *gin.Context) {
	var requestBody api.CreateArticleJSONRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(400, api.ErrorResponse{Message: err.Error()})
		return
	}
}
