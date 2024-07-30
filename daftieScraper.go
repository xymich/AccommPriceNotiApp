package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

// order of callback calls :=
// 1. OnRequest(): Called before performing an HTTP request with Visit().
// 2. OnError(): Called if an error occurred during the HTTP request.
// 3. OnResponse(): Called after receiving a response from the server.
// 4. OnHTML(): Called right after OnResponse() if the received content is HTML.
// 5. OnScraped(): Called after all OnHTML() callback executions are completed.

func main() {

	scrapeUrl := "https://go-colly.org/docs/introduction/configuration/"

	clt := colly.NewCollector(colly.AllowedDomains("www.daft.ie", "daft.ie", "google.com", "go-colly.org"))

	// clt.OnHTML("h2.collector-configuration", func(h *colly.HTMLElement) {
	// 	fmt.Println(h.Text)
	// })

	clt.OnHTML(".language-go", func(e *colly.HTMLElement) {
		//fmt.Println(e)
		fmt.Println(e.ChildText("span"))
	})

	clt.OnRequest(func(req *colly.Request) {
		fmt.Println("Visiting: ", req.URL)
	})

	clt.OnError(func(req *colly.Response, err error) {
		fmt.Printf("Error while scraping: %s\n", err.Error())
	})

	clt.Visit(scrapeUrl)

	fmt.Println()

}
