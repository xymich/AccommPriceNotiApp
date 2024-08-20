package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

// https://www.daft.ie/property-for-sale/XXXX?pageSize=20&from=0 // > can go to max on bottom
var TempForSaleUrl = "https://www.daft.ie/property-for-sale/dundalk-louth"
var TempForRentUrl = "https://www.daft.ie/property-for-rent/monaghan"
var PlaywrightContext playwright.BrowserContext

type DaftComponents struct {
	// AdvertLink		string 	`json:"advert_link"`
	Address			string	`json:"address"`
	Price			string 	`json:"price"`
	BedCount		string 	`json:"bed_count"`
	BathCount		string 	`json:"bath_count"`
	Size			string 	`json:"size"`
	PropertyType	string 	`json:"property_type"`
	Seller			string 	`json:"seller"`
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
	runtimes1 := time.Now()


	compArr := PageScrape(PlaywrightContext)
	fmt.Print(compArr)

	runtimes2 := time.Now()
	fmt.Println(runtimes2.Sub(runtimes1).Seconds)
	PlaywrightContext.Close()
}

func PageScrape(ctx playwright.BrowserContext) (data []DaftComponents) {
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

	// screenshot, err := page.Screenshot(playwright.PageScreenshotOptions{
	// 	Path: playwright.String("helloworld.png"),
	// })

	if err != nil {
		log.Fatalf("could not screenshot: %v", err)
	}

	lixPath := fmt.Sprintf("xpath=//html/body/div[2]/main/div[3]/div[1]/ul/li") //reformat as you wish
	liLocators, err := page.Locator(lixPath).All()
	if err != nil {
		log.Fatalf("could not get entries: %v", err)
	}

	if len(liLocators) == 0 {
		log.Fatalf("no entries found")
	}

	dataEntries := []DaftComponents{}

	for r:=0;r<len(liLocators)-1;r++ {
		
		liInnerTexts, err := liLocators[r].AllInnerTexts()
		if err != nil {
			log.Fatalf("could not get entries: %v", err)
		}

		liInnerSplit := strings.Split(liInnerTexts[0],"\n")

		houseDataEntry := createDataEntry(liInnerSplit)

	dataEntries = append(dataEntries, houseDataEntry)

	for _,v:=range liInnerSplit{
		fmt.Println(v)
	}
	
	fmt.Print("\n\n)))(((**********************)))(((\n\n")
	}

	fmt.Print(dataEntries[0], "\n\n******************n\n")

	return dataEntries
}

func createDataEntry(liInnerSplit []string) (dataEntry DaftComponents) {

	// dataEntry = DaftComponents{
	// 	liInnerSplit[0],
	// 	liInnerSplit[1],
	// 	liInnerSplit[2],
	// 	liInnerSplit[3],
	// 	liInnerSplit[4],
	// 	liInnerSplit[5],
	// 	liInnerSplit[6],
	// }
	
	return dataEntry
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