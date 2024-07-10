package crawler

import (
	"log"
	"time"

	"github.com/gocolly/colly"
)

type LinkExtractor struct {
	c         *colly.Collector
	urlString string
}

func NewExtractor(collector *colly.Collector, scrape_url string) *LinkExtractor {
	return &LinkExtractor{
		c:         collector,
		urlString: scrape_url,
	}
}

// get all TitleHeadings and send to Kafka topic
func (s *LinkExtractor) ExtractArticles(linkHandler func(string)) {
	// config stuff so we get authorized
	s.c.AllowURLRevisit = true

	s.c.Limit(&colly.LimitRule{
		RandomDelay: time.Millisecond * 200,
	})

	s.c.OnError(func(_ *colly.Response, err error) {
		log.Println("Error: ", err)
	})

	s.c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.5")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("DNT", "1")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
	})

	// html callbacks

	s.c.OnHTML("a[data-testid='TitleLink']", func(e *colly.HTMLElement) {
		if article_link := e.Attr("href"); article_link != "" {
			log.Printf("Got link: %s", article_link)
			linkHandler(article_link)
		} else {
			log.Println("Error finding article href")
		}
	})

	s.c.OnHTML("a[data-testid='Title']", func(e *colly.HTMLElement) {
		if article_link := e.Attr("href"); article_link != "" {
			log.Printf("Got link: %s", article_link)
			linkHandler(article_link)
		} else {
			log.Println("Error finding article href")
		}
	})

	s.c.Visit(s.urlString)
}
