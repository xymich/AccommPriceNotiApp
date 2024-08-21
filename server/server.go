package server

import (
	"github.com/ilyTea/AccommPriceNotiApp/scraper"
)

func Start() {
	scraper.Scrape()
}