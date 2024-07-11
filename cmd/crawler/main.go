package main

import (
	"log"
	"news-crawler/internal/crawler"
	"news-crawler/internal/queue"
)

func main() {
	rabbitMQ, err := queue.NewRabbitMQ("amqp://guest:guest@localhost:5672/", "article_links")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	urlString := "https://www.reuters.com/"
	linkExtractor := crawler.NewExtractor(urlString)

	linkHandler := func(link string) {
		err := rabbitMQ.PublishMessage(link)
		if err != nil {
			log.Printf("Error publishing ot RabbitMQ: %v", err)
		}
	}

	linkExtractor.ExtractArticles(linkHandler)
}
