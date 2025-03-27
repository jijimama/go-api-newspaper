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

// 構造体をjsonに変換する
func (a *Article) MarshalJSON() ([]byte, error) {
	return json.Marshal(&api.ArticleResponse{
		Id:          a.ID,
		Body:        a.Body,
		Year:        a.Year,
		Month:       a.Month,
		Day:         a.Day,
		Newspaper:   api.Newspaper{
			Id:   &a.Newspaper.ID,
			Name: a.Newspaper.Name,
		}
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
		Newspaper:   newspaper,
		NewspaperID: newspaper.ID,
	}
	if err := DB.Create(article).Error; err != nil {
		return nil, err
	}
	return article, nil
}
