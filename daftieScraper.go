package main

import (
	"fmt"
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

	compArr := ExtractComponents(PlaywrightContext)
	fmt.Print(compArr)
	PlaywrightContext.Close()
}

func ExtractComponents(ctx playwright.BrowserContext) (data []DaftComponents) {

	// Created a new page from the context we initialized
	page, err := ctx.NewPage()

	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	// Navigates to the Corporate Announcements Page
	if _, err = page.Goto(TempForSaleUrl, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	// Waits until the full URL is loaded
	err = page.WaitForURL(TempForSaleUrl)

	if err != nil {
		log.Fatalf("could not wait for url: %v", err)
	}

	screenshot, err := page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("foo.png"),
	})
	screenshot = screenshot

	if err != nil {
		log.Fatalf("could not screenshot: %v", err)
	}

	listItem := 2
	_xpath := fmt.Sprintf("xpath=//html/body/div[2]/main/div[3]/div[1]/ul/li[%v]", listItem)
  sampleHeading, err := page.Locator(_xpath).AllTextContents()
	if err != nil {
		log.Fatalf("could not get entries: %v", err)
	}

	if len(sampleHeading) == 0 {
		log.Fatalf("no entries found")
	}
	fmt.Printf("%v",sampleHeading)
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
		UserAgent: playwright.String("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36"), // Fake user agent
		Viewport:  &playwright.Size{Width: 1920, Height: 1080},
	})
	if err != nil {
		log.Fatalf("could not create context: %v", err)
	}

	fmt.Printf("\n\nwowwww\n\n\n")



	return context

}