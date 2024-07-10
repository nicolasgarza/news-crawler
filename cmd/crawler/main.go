package main

import (
	"log"
	"news-crawler/internal/crawler"
	"news-crawler/internal/queue"

	"github.com/gocolly/colly"
)

func main() {
	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	rabbitMQ, err := queue.NewRabbitMQ("amqp://guest:guest@localhost:5672/", "article_links")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	urlString := "https://www.reuters.com/"
	linkExtractor := crawler.NewExtractor(collector, urlString)

	linkHandler := func(link string) {
		err := rabbitMQ.PublishMessage(link)
		if err != nil {
			log.Printf("Error publishing ot RabbitMQ: %v", err)
		}
	}

	linkExtractor.ExtractArticles(linkHandler)
}
