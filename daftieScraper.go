package main

import (
	"log"

	"github.com/playwright-community/playwright-go"
)

var PlaywrightContext playwright.BrowserContext

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

func main() {
	err := playwright.Install()
	if err != nil {
		log.Fatalf("could not install playwright: %v", err)
	}

	PlaywrightContext = InitializePlaywright()

}
