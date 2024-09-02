package scraper

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
)

var propertyType string = "rent"
var PlaywrightContext playwright.BrowserContext

type DaftComponents struct {
	Address			string	`json:"address"`
	Price			string 	`json:"price"`
	BedCount		string 	`json:"bed_count"`
	BathCount		string 	`json:"bath_count"`
	Size			string 	`json:"size"`
	PropertyType	string 	`json:"property_type"`
	Seller			string 	`json:"seller"`
	//AdvertLink	string 	`json:"advert_link"` // ?? maybe
	//PropertyImage	string 	`json:"property_image"` // ?? maybe
}

func Scrape() {
	// extra check specifically for server use (might be needed just incase but also might be worth removing if works without to increase speed)
	err := playwright.Install()
	if err != nil {
		log.Fatalf("could not install playwright: %v", err)
	}
	PlaywrightContext = InitializePlaywright()

	//now we scrape!
	scrapeUrl := fmt.Sprintf("https://www.daft.ie/property-for-%v/ireland", propertyType)
	listCount := 0
	compArr, totalListing := pageScrape(scrapeUrl, PlaywrightContext)
	fmt.Println("scrapeurl = ", scrapeUrl)
	pageCount:= totalListing/20;
	allData := [][]DaftComponents{compArr}
	dataChan := make(chan [][]DaftComponents)
	
	pageScrapeRoutine := func( init int, limit int, wg *sync.WaitGroup) {
		fmt.Println("time before context",time.Now())
		defer wg.Done()
		var pctx = InitializePlaywright()
		scrapedGroupData := pageScrapeIncrement(pctx, init, limit)
		dataChan <- scrapedGroupData
	}

	if (pageCount > 1) {
		numOfRoutines := 8
		pageRemainder := pageCount % numOfRoutines
		pagesPerRoutine := (pageCount - pageRemainder) / numOfRoutines
		initPage := 1
		limitPage := pagesPerRoutine
		
		var limits [8][2]int
		for i := 0; i < numOfRoutines; i++ {
			limits[i][0] = initPage
			limits[i][1] = limitPage
			initPage += pagesPerRoutine
			limitPage += pagesPerRoutine
		}
		fmt.Println("limits arr ",limits)
		go func() {
			wg := sync.WaitGroup{}
			for i := 0; i < numOfRoutines; i++ {
				wg.Add(1)
				go pageScrapeRoutine(limits[i][0], limits[i][1], &wg)
			}
			wg.Wait()
			close(dataChan)
		}()
	fmt.Println(len(dataChan))
	  for n:= range dataChan {
		allData = append(allData, n...) 
	  }
	}
	fmt.Println("final listcount :", listCount,"\n\n >>>>>><<<<<<")
	PlaywrightContext.Close()
}

func pageScrapeIncrement(ctx playwright.BrowserContext,initial int, limit int) (data [][]DaftComponents ) {
	// https://www.daft.ie/property-for-rent/ireland?from=20&pageSize=20
	groupedData := [][]DaftComponents{}

	for initial < limit {
		initial += 1
		url := fmt.Sprintf("https://www.daft.ie/property-for-rent/ireland?from=%v&pageSize=20",initial*20)
		daftComponents, _ := pageScrape(url, ctx)
		groupedData = append(groupedData,daftComponents)
	}

	return groupedData
}

func pageScrape(url string, ctx playwright.BrowserContext) (data []DaftComponents, totalPageCount int ) {
	// Created a new page from the context we initialized
	page, err := ctx.NewPage()

	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	// Navigates to the Corporate Announcements Page
	if _, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	// Waits until the full URL is loaded
	err = page.WaitForURL(url)
	if (err != nil) {
		fmt.Println(err)
	}
	screenshot, err := page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("helloworld.png"),
	})
	screenshot = screenshot

	if err != nil {
		log.Fatalf("could not wait for url: %v", err)
	}

	if err != nil {
		log.Fatalf("could not screenshot: %v", err)
	}

	lixPath := "xpath=//html/body/div[2]/main/div[3]/div[1]/ul/li" //reformat as you wish
	liLocators, err := page.Locator(lixPath).All()
	if err != nil {
		log.Fatalf("could not get entries: %v", err)
	}

	if len(liLocators) == 0 {
		log.Fatalf("no entries found")
	}

	dataEntries := []DaftComponents{}

	fmt.Println(len(liLocators))

	for r:=0;r<len(liLocators)-1;r++ {
		
		liInnerTexts, err := liLocators[r].AllInnerTexts()
		if err != nil {
			log.Fatalf("could not get entries: %v", err)
		}

		liInnerSplit := strings.Split(liInnerTexts[0],"\n")

		houseDataEntry := createDataEntry(liInnerSplit)

		dataEntries = append(dataEntries, houseDataEntry)

	}
	
	paginationTextArray, _ := page.Locator("xpath=//html/body/div[2]/main/div[3]/div[1]/div[2]/p").AllInnerTexts()
	splitPaginationText := strings.Split(paginationTextArray[0], " ")

	totalListCount, err := strconv.Atoi(strings.Replace(splitPaginationText[len(splitPaginationText)-1], ",", "", -1))
	if err != nil {
		log.Fatalf("Could not get total list count: %v", err)
	}

	return dataEntries, totalListCount
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
	fmt.Println(liInnerSplit, "\n\n****************************************\n")
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