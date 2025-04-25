package models

import (
	"encoding/json"
	"go-api-newspaper/api"
)

type Article struct {
	ID          int
	Body        string
	Year        int
	Month       int
	Day         int
	NewspaperID int
	Newspaper   *Newspaper
}

func (a *Article) MarshalJSON() ([]byte, error) {
	return json.Marshal(&api.ArticleResponse{
		Id:    a.ID,
		Body:  a.Body,
		Year:  a.Year,
		Month: a.Month,
		Day:   a.Day,
		Newspaper: api.Newspaper{
			Id:         &a.Newspaper.ID,
			Title:      a.Newspaper.Title,
			ColumnName: a.Newspaper.ColumnName,
		},
	})
}

func CreateArticle(body string, year int, month int, day int, newspaperID int) (*Article, error) {
	newspaper, err := GetNewspaper(newspaperID)
	if err != nil {
		return nil, err
	}

	article := &Article{
		Body:        body,
		Year:        year,
		Month:       month,
		Day:         day,
		NewspaperID: newspaperID,
		Newspaper:   newspaper,
	}
	if err := DB.Create(article).Error; err != nil {
		return nil, err
	}
	return article, nil
}

func GetArticle(id int) (*Article, error) {
	article := &Article{}
	if err := DB.Where("id = ?", id).First(article).Error; err != nil {
		return nil, err
	}
	return article, nil
}
