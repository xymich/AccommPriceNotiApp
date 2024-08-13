package main

import (
	"log"

	"AccomPriceNotiApp/services"

	"github.com/playwright-community/playwright-go"
)

func main() {
	err := playwright.Install()
	if err != nil {
		log.Fatalf("could not install playwright: %v", err)
	}

	services.PlaywrightContext = services.InitializePlaywright()
}
