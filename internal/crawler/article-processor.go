package crawler

import (
	"log"
	"news-crawler/internal/utils"
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

func NewScraper(articleUrl string) *ArticleScraper {
	return &ArticleScraper{
		c:          utils.NewCollector(),
		articleUrl: articleUrl,
	}
}

func (a *ArticleScraper) ScrapeArticle() (*articleInfo, error) {
	articleInfo := &articleInfo{}
	var scrapeErr error

	a.c.OnError(func(_ *colly.Response, err error) {
		scrapeErr = err
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

	err := a.c.Visit(a.articleUrl)
	if err != nil {
		return nil, err
	}

	if scrapeErr != nil {
		return nil, scrapeErr
	}

	return articleInfo, nil
}

func parseDate(dateStr string) (time.Time, error) {
	layout := "January 2, 2006 3:04 PM MST"
	return time.Parse(layout, dateStr)
}
