package crawlers

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

func BeginBBCScrape() {
	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		//ensure internal links are followed

		// if strings.Contains(e.Attr("href"), "https://www.bbc.co.uk/sport/football") {

		if strings.Contains(e.Attr("href"), "colly") {
			e.Request.Visit(e.Attr("href"))
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://go-colly.org/")
}
