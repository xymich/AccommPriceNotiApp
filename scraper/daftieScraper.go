package scraper

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/playwright-community/playwright-go"
)

const urlPre string = "https://www.daft.ie/property-for-rent/"
const urlExt string = "?pageSize=20&from="
var location string = ""

//@@// testPageUrl = urlPre + location + urlExt + (listCount+=20) ie. https://www.daft.ie/property-for-sale/dundalk-louth?pageSize=20&from=0

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

func Scrape(location string) {
	// extra check specifically for server use (might be needed just incase but also might be worth removing if works without to increase speed)
	err := playwright.Install()
	if err != nil {
		log.Fatalf("could not install playwright: %v", err)
	}
	PlaywrightContext = InitializePlaywright()

	//now we scrape!
	scrapeUrl := ""
	listCount := 0

	for {
		goNextPage := false
		scrapeUrl = fmt.Sprintf(urlPre + location + urlExt + strconv.Itoa(listCount))
		compArr, goNextPage := pageScrape(scrapeUrl, PlaywrightContext)

		listCount += 20

		compArr = compArr
		fmt.Println("go next page? :", goNextPage, "|| listcount :", listCount,"\n\n >>>>>><<<<<<")
		if (goNextPage == false) {
			break
		}
	}
	
	PlaywrightContext.Close()
}

func pageScrape(url string, ctx playwright.BrowserContext) (data []DaftComponents, goNextPage bool ) {
	// Created a new page from the context we initialized
	goNextPage = false
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

	// screenshot, err := page.Screenshot(playwright.PageScreenshotOptions{
	// 	Path: playwright.String("helloworld.png"),
	// })

	if err != nil {
		log.Fatalf("could not wait for url: %v", err)
	}

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
	
	fmt.Print("\n]][[###**********^v^**********###]][[\n\n")
	}
	
	paginationTextArray, err := page.Locator("xpath=//html/body/div[2]/main/div[3]/div[1]/div[2]/p").AllInnerTexts()
	splitPaginationText := strings.Split(paginationTextArray[0], " ")
	currentListCount,err := strconv.Atoi(strings.Replace(splitPaginationText[len(splitPaginationText)-3], ",", "", -1))
	if err != nil {
		log.Fatalf("Could not get current list count: %v", err)
	}

	totalListCount,err := strconv.Atoi(strings.Replace(splitPaginationText[len(splitPaginationText)-1], ",", "", -1))
	if err != nil {
		log.Fatalf("Could not get total list count: %v", err)
	}

	fmt.Println("CURRENT LIST COUNT ",currentListCount)
	fmt.Println("TOTAL LIST COUNT ",totalListCount)

	if (totalListCount > currentListCount) {
		goNextPage = true;
	}

	return dataEntries, goNextPage
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