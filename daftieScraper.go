package main

import (
	"log"

	"github.com/playwright-community/playwright-go"
)

var TempForSaleUrl = "https://www.daft.ie/property-for-sale/monaghan"
var TempForRentUrl = "https://www.daft.ie/property-for-rent/monaghan"
var PlaywrightContext playwright.BrowserContext

type DaftComponents struct {
	AdvertLink   string `json:"advert_link"`
	PaymentType  string `json:"payment_type"`
	PropertyType string `json:"property_type"`
	Title        string `json:"title"`
	Price        string `json:"price"`
	BedCount     string `json:"bed_count"`
	Tag          string `json:"tag"`
	//PropertyImage	string `json:"property_image"` // ?? maybe
}

func main() {
	// extra check specifically for server use
	err := playwright.Install()
	if err != nil {
		log.Fatalf("could not install playwright: %v", err)
	}
	//init
	PlaywrightContext = InitializePlaywright()
	//now we scrape!

	PlaywrightContext.Close()
}

func ExtractDaftComponents(ctx playwright.BrowserContext) (data []DaftComponents) {

	// Created a new page from the context we initialized
	page, err := ctx.NewPage()

	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	// Navigates to the Corporate Announcements Page
	if _, err = page.Goto("https://www.bseindia.com/corporates/ann.html", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	// Waits until the full URL is loaded
	err = page.WaitForURL("https://www.bseindia.com/corporates/ann.html")

	if err != nil {
		log.Fatalf("could not wait for url: %v", err)
	}

	return
}

func InitializePlaywright() playwright.BrowserContext {
	// Installation of browser and OS dependencies
	err := playwright.Install()
	if err != nil {
		log.Fatalf("could not install playwright: %v", err)
	}

	// Running playwright
	pw, err := playwright.Run(
		&playwright.RunOptions{},
	)
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}

	// Selecting Chromium as the browser
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Args: []string{"--disable-blink-features=AutomationControlled"},
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

	// Creating context out of the created browser
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		UserAgent: playwright.String("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"), // Fake user agent
		Viewport:  &playwright.Size{Width: 1920, Height: 1080},
	})
	if err != nil {
		log.Fatalf("could not create context: %v", err)
	}

	return context

}
