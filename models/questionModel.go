package models

type Question struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Location string `json:"location"`
}
