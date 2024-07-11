package models

import "time"

type Metadata struct {
	ID        string    `bson:"_id,omitempty"`
	ArticleID string    `bson:"article_id"`
	URL       string    `bson:"url"`
	ScrapedAt time.Time `bson:"scraped_at"`
	Source    string    `bson:"source"`
}
