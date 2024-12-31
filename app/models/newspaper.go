package models

import (
	"encoding/json"

	"go-api-newspaper/api"
)

type Newspaper struct {
	ID          int
	Title       string
	ColumnName  string
}

// 構造体をjsonに変換する
func (a *Newspaper) MarshalJSON() ([]byte, error) {
	return json.Marshal(&api.NewspaperResponse{ // api.NewspaperResponse という別の構造体にデータを詰め替えている
		Id:          a.ID,
		Title:       a.Title,
		ColumnName:  a.ColumnName,
	})
}

func CreateNewspaper(title string, columnName string) (*Newspaper, error) {
	newspaper := &Newspaper{
		Title:       title,
		ColumnName:  columnName,
	}
	if err := DB.Create(newspaper).Error; err != nil {
		return nil, err
	}
	return newspaper, nil
}

func GetNewspaper(ID int) (*Newspaper, error) {
	var newspaper = Newspaper{}
	if err := DB.First(&newspaper, ID).Error; err != nil {
		return nil, err
	}
	return &newspaper, nil
}

func (a *Newspaper) Save() error {
	if err := DB.Save(&a).Error; err != nil {
		return err
	}
	return nil
}

func (a *Newspaper) Delete() error {
	if err := DB.Where("id = ?", &a.ID).Delete(&a).Error; err != nil {
		return err
	}
	return nil
}
