package main

import (
	"log"
	"news-crawler/internal/crawler"

	"github.com/gocolly/colly"
)

func main() {
	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)
	articleString := "https://www.reuters.com/markets/deals/mountain-asset-sales-loom-after-oil-megamerger-era-2024-06-26/"

	articleProcessor := crawler.NewScraper(collector, articleString)
	articleInfo, err := articleProcessor.ScrapeArticle()
	if err != nil {
		log.Printf("Failed to scrape articles: %v", err)
	}

	log.Printf("Article title: %s", articleInfo.Title)
	log.Printf("Article authors: %q\n", &articleInfo.Authors)
	log.Printf("Article date: %v\n", &articleInfo.Date)
	log.Printf("Article content length: %d characters", len(articleInfo.Content))
}
