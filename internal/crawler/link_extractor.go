package crawler

import (
	"log"
	"net/url"
	"news-crawler/internal/utils"
	"strings"

	"github.com/gocolly/colly"
)

type LinkExtractor struct {
	c         *colly.Collector
	urlString string
	baseURL   string
}

func NewExtractor(scrape_url string) *LinkExtractor {
	parsedURL, err := url.Parse(scrape_url)
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
	}
	baseURL := parsedURL.Scheme + "://" + parsedURL.Host
	return &LinkExtractor{
		c:         utils.NewCollector(),
		urlString: scrape_url,
		baseURL:   baseURL,
	}
}

// get all TitleHeadings and send to Kafka topic
func (s *LinkExtractor) ExtractArticles(linkHandler func(string)) {
	// html callbacks

	handleLink := func(e *colly.HTMLElement) {
		if article_link := e.Attr("href"); article_link != "" {
			fullURL := s.getFullURL(article_link)
			log.Printf("Got link: %s", fullURL)
			linkHandler(fullURL)
		} else {
			log.Println("Error finding article href")
		}
	}

	s.c.OnError(func(_ *colly.Response, err error) {
		log.Println("Error: ", err)
	})

	s.c.OnHTML("a[data-testid='TitleLink']", handleLink)

	s.c.OnHTML("a[data-testid='Title']", handleLink)

	s.c.Visit(s.urlString)
}

func (s *LinkExtractor) getFullURL(articleLink string) string {
	if strings.HasPrefix(articleLink, "http://") || strings.HasPrefix(articleLink, "https://") {
		return articleLink
	}
	return s.baseURL + articleLink
}
