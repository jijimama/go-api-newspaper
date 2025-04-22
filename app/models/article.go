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
