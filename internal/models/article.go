package models

import "time"

type Article struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Authors []string  `json:"authors"`
	Date    time.Time `json:"date"`
	Content string    `json:"content"`
	URL     string    `json:"url"`
}
