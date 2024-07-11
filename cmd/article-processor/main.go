package main

import (
	"log"
	"news-crawler/internal/crawler"
	"news-crawler/internal/queue"
	"time"
)

func main() {
	rabbitMQ, err := queue.NewRabbitMQ("amqp://guest:guest@localhost:5672/", "article_links")
	if err != nil {
	}

	msgs, err := rabbitMQ.Consume()
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("Waiting for messages, ctrl+c to exit")

	for msg := range msgs {
		articleUrl := string(msg.Body)
		log.Printf("Received a message: %s", articleUrl)

		articleProcessor := crawler.NewScraper(articleUrl)
		articleInfo, err := articleProcessor.ScrapeArticle()
		if err != nil {
			log.Printf("Failed to scrape article %s: %v", articleUrl, err)
		}

		log.Printf("Article title: %s", articleInfo.Title)
		log.Printf("Article authors: %q", &articleInfo.Authors)
		log.Printf("Article date: %v", &articleInfo.Date)
		log.Printf("Article content length: %d characters", len(articleInfo.Content))
		log.Println("-----------------------------------")
		time.Sleep(time.Second * 3)

	}
}
