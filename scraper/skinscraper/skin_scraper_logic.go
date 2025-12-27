package skinsscraper

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/eovacius/csgodatabase-scraper/internal"
	"github.com/eovacius/csgodatabase-scraper/scraper"
	"github.com/eovacius/csgodatabase-scraper/scraper/config"
)

func ScrapeList(ctx context.Context, list []string, key string) []config.Skin {
	var allSkins []config.Skin
	path := "json/" + key
	_ = os.Mkdir(path, 0755)

	const maxRetries = 2

	for _, item := range list {
		url := config.Target + "/" + key + "/" + item + "/"

		var skins []config.Skin
		success := false

		for attempt := 0; attempt <= maxRetries && !success; attempt++ {
			if attempt == 0 {
				fmt.Printf("Scraping: %s\n", url)
			} else {
				fmt.Printf("Retry %d/%d: %s\n", attempt, maxRetries, url)
			}

			var pageTitle string
			var rawData []map[string]string

			err := chromedp.Run(ctx,
				chromedp.Navigate(url),
				chromedp.Evaluate(string(scraper.ConfigJS), nil),
				chromedp.Sleep(config.Delay),
				chromedp.Title(&pageTitle),
				chromedp.Evaluate(string(scraper.ScriptJS), &rawData),
			)

			if err != nil {
				fmt.Printf("\033[31m[!]\033[0m chromedp error: %v\n", err)
				break
			}

			if strings.Contains(strings.ToLower(pageTitle), "page not found") {
				fmt.Printf("\033[33m[?]\033[0m 404 Page not found: %s\n", url)
				break
			}

			if strings.Contains(strings.ToLower(pageTitle), "verify") ||
				strings.Contains(strings.ToLower(pageTitle), "human") ||

				// thats to test detection retry logic
				// strings.Contains(strings.ToLower(pageTitle), "the 2018 inferno collection skins - csgo database") ||
				strings.Contains(strings.ToLower(pageTitle), "just a moment") {

				if attempt == 0 {
					fmt.Printf("\033[31m[!]\033[0m Detection triggered! Switching to stealth mode...\n")
				} else {
					fmt.Printf("\033[31m[!]\033[0m Detection triggered!\n")
				}
				randomMs := 500 + rand.Intn(1500)
				config.Delay += time.Duration(randomMs) * time.Millisecond
				continue
			}

			for _, data := range rawData {
				skins = append(skins, config.Skin{
					Name:       data["name"],
					Weapon:     internal.SpecialMark(data["weapon"]),
					Rarity:     data["rarity"],
					Collection: data["collection"],
					Price:      internal.ParsePrice(data["price"], data["stattrakPrice"]),
					URL:        data["url"],
				})
			}

			if len(skins) == 0 {
				fmt.Printf("\033[31m[!]\033[0m No skins found\n")
				continue
			}
			success = true
		}

		if !success {
			fmt.Printf("\033[31m[!]\033[0m Failed to scrape: %s\n", url)
			continue
		}

		filename := fmt.Sprintf("%s/%s.json", path, item)
		internal.SaveJSON(filename, skins)
		allSkins = append(allSkins, skins...)
	}

	return allSkins
}
