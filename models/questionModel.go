package models

import "time"

type Question struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	CreatedAt time.Time     `json:"created_at"`
	Content   string        `json:"content"`
	Location  time.Location `json:"location"`
}

type QuestionRequest struct {
	Title    string        `json:"title"`
	Content  string        `json:"content"`
	Location time.Location `json:"location"`
}

type QuestionResponse struct {
	Title     string        `json:"title"`
	CreatedAt time.Time     `json:"created_at"`
	Content   string        `json:"content"`
	Location  time.Location `json:"location"`
}

type SearchResponse struct {
	Time      string             `json:"time"`
	Hits      string             `json:"hits"`
	Questions []QuestionResponse `json:"documents"`
}
