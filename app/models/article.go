package models

type Article struct {
	ID          int
	Body        string
	Year        int
	Month       int
	Day         int
	NewspaperID int
	Newspaper   *Newspaper
}
