package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {

	scrapeUrl := "https://www.daft.ie"

	clt := colly.NewCollector(colly.AllowedDomains("www.daft.ie", "daft.ie"))

	clt.OnHTML("h1.homepage-tagline", func(header *colly.HTMLElement) {
		fmt.Println(header.Text)
	})

	clt.OnRequest(func(req *colly.Request) {
		fmt.Printf(fmt.Sprintf("Visiting %s\n", req.URL))
	})

	clt.OnError(func(req *colly.Response, err error) {
		fmt.Printf("Error while scraping: %s\n", err.Error())
	})

	clt.Visit(scrapeUrl)
	fmt.Println(222)
}
