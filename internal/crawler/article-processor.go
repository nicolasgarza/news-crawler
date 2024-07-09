package crawler

import (
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type ArticleScraper struct {
	c          *colly.Collector
	articleUrl string
}

type articleInfo struct {
	Title   string
	Authors []string
	Date    time.Time
	Content string
}

func NewScraper(collector *colly.Collector, articleUrl string) *ArticleScraper {
	return &ArticleScraper{
		c:          collector,
		articleUrl: articleUrl,
	}
}

func (a *ArticleScraper) ScrapeArticle() {
	articleInfo := &articleInfo{}

	// config stuff so we get authorized
	a.c.AllowURLRevisit = true

	a.c.Limit(&colly.LimitRule{
		RandomDelay: time.Millisecond * 200,
	})

	a.c.OnError(func(_ *colly.Response, err error) {
		log.Println("Error: ", err)
	})

	a.c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.5")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("DNT", "1")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
	})

	// get title
	a.c.OnHTML("div[class='default-article-header__heading__3cyKI']", func(e *colly.HTMLElement) {
		headerText := e.ChildText("h1")

		if headerText != "" {
			articleInfo.Title = headerText
		} else {
			log.Println("Error getting header text")
		}
	})

	// get date
	a.c.OnHTML("time[data-testid='Body']", func(e *colly.HTMLElement) {
		var dateTexts []string
		e.ForEach("span[class='date-line__date___kNbY']", func(i int, el *colly.HTMLElement) {
			if i < 2 {
				dateTexts = append(dateTexts, el.Text)
			}
		})

		fullDateText := strings.Join(dateTexts, " ")
		parsedDate, err := parseDate(fullDateText)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
		} else {
			articleInfo.Date = parsedDate
		}
	})

	// get author
	a.c.OnHTML("div[class='info-content__author-date__1Epi_']", func(e *colly.HTMLElement) {
		e.ForEach("a[rel='author']", func(_ int, el *colly.HTMLElement) {
			articleInfo.Authors = append(articleInfo.Authors, el.Text)
		})

		if len(articleInfo.Authors) == 0 {
			log.Println("No authors found")
		}
	})

	// get article text
	var articleContent strings.Builder
	a.c.OnHTML("div[data-testid^='paragraph-']", func(e *colly.HTMLElement) {
		paragraphText := e.Text
		articleContent.WriteString(paragraphText)
		articleContent.WriteString("\n\n")
	})

	a.c.OnScraped(func(r *colly.Response) {
		articleInfo.Content = articleContent.String()
	})

	a.c.Visit(a.articleUrl)
	log.Printf("Article title: %s", articleInfo.Title)
	log.Printf("Article authors: %q\n", &articleInfo.Authors)
	log.Printf("Article date: %v\n", &articleInfo.Date)
	log.Printf("Article content length: %d characters", len(articleInfo.Content))
}

func parseDate(dateStr string) (time.Time, error) {
	layout := "January 2, 2006 3:04 PM MST"
	return time.Parse(layout, dateStr)
}
