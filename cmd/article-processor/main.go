package main

import (
	"log"
	"news-crawler/internal/crawler"
	"news-crawler/internal/models"
	"news-crawler/internal/queue"
	"news-crawler/internal/storage"
	"news-crawler/internal/utils"
	"news-crawler/internal/workers"
	"time"
)

func main() {
	rabbitMQ, err := queue.NewRabbitMQ("amqp://guest:guest@localhost:5672/", "article_links")
	if err != nil {
	}

	esClient, err := storage.NewElasticSearchClient([]string{"http://localhost:9200"}, "articles")
	if err != nil {
		log.Fatalf("Failed to connect to Elasticsearch: %v", err)
	}

	mongoClient, err := storage.NewMongoDBClient("mongodb://localhost:27017", "news_aggregator", "metadata")
	if err != nil {
		log.Fatalf("Failed to connect to Mongodb: %v", err)
	}

	workerPool := workers.NewWorkerPool(5)
	workerPool.Start()
	defer workerPool.Stop()

	msgs, err := rabbitMQ.Consume()
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("Waiting for messages, ctrl+c to exit")

	for msg := range msgs {
		articleUrl := string(msg.Body)
		workerPool.Submit(func() {
			processArticle(articleUrl, esClient, mongoClient)
		})
	}
}

func processArticle(articleUrl string, esClient *storage.ElasticSearchClient, mongoClient *storage.MongoDBClient) {
	log.Printf("Received a message: %s", articleUrl)

	articleProcessor := crawler.NewScraper(articleUrl)
	articleInfo, err := articleProcessor.ScrapeArticle()
	if err != nil {
		log.Printf("Failed to scrape article %s: %v", articleUrl, err)
		log.Println("-----------------------------------")
		return
	}

	article := &models.Article{
		ID:      utils.GenerateID(),
		Title:   articleInfo.Title,
		Authors: articleInfo.Authors,
		Date:    articleInfo.Date,
		Content: articleInfo.Content,
		URL:     articleUrl,
	}

	err = esClient.IndexArticle(article)
	if err != nil {
		log.Printf("Failed to index article in Elasticsearch: %v", err)
	}

	metadata := &models.Metadata{
		ArticleID: article.ID,
		URL:       articleUrl,
		ScrapedAt: time.Now(),
		Source:    "Reuters",
	}

	err = mongoClient.InsertMetadata(metadata)
	if err != nil {
		log.Printf("Failed to insert metadata in MongoDB: %v", err)
	}

	log.Println("Processed article: %s", article.Title)
	time.Sleep(time.Second * 1)
}
