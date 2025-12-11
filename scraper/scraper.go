package scraper

import (
	"context"
	_ "embed"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/eovacius/csgodatabase-scraper/internal"
	"github.com/eovacius/csgodatabase-scraper/scraper/config"
)

//go:embed config/scripts/script.js
var ScriptJS string

//go:embed config/scripts/config.js
var ConfigJS string

// Run func executes the scraping process and returns a slice of Skin structs with error if any.
func ScrapeSkins() ([]config.Skin, error) {
	// fancy ascii title ;)
	fmt.Print("\033[34m" + `
   ____  ____     ____    ____      ____    ____        _       ____   U _____ u   ____     
U /"___|/ __"| u |___"\  / __"| uU /"___|U |  _"\ u U  /"\  u U|  _"\ u\| ___"|/U |  _"\ u  
\| | u <\___ \/  U __) |<\___ \/ \| | u   \| |_) |/  \/ _ \/  \| |_) |/ |  _|"   \| |_) |/  
 | |/__ u___) |  \/ __/ \u___) |  | |/__   |  _ <    / ___ \   |  __/   | |___    |  _ <    
  \____||____/>> |_____|u|____/>>  \____|  |_| \_\  /_/   \_\  |_|      |_____|   |_| \_\   
 _// \\  )(  (__)<<  //   )(  (__)_// \\   //   \\_  \\    >>  ||>>_    <<   >>   //   \\_  
(__)(__)(__)    (__)(__) (__)    (__)(__) (__)  (__)(__)  (__)(__)__)  (__) (__) (__)  (__) 
		` + "\033[0m")

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), config.Opts...)
	defer cancel()

	fmt.Println("\n[*] Creating allocator context and applying opts...")

	browserCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	browserCtx, cancel = context.WithTimeout(browserCtx, config.DeadLine)
	defer cancel()

	fmt.Println("[*] Created allocator")

	var wg sync.WaitGroup
	var cases, collections []config.Skin

	spawnWorker(browserCtx, &wg, &cases, config.List, "cases")
	spawnWorker(browserCtx, &wg, &collections, config.CollectionsList, "collections")

	// thats to test whether we get detected or not
	// spawnWorker(browserCtx, &wg, &cases, config.List, "cases")
	// spawnWorker(browserCtx, &wg, &collections, config.CollectionsList, "collections")

	wg.Wait()

	merged := append(cases, collections...)

	uniqueSkins := internal.RemoveDuplicates(merged)
	return uniqueSkins, nil
}

func spawnWorker(browserCtx context.Context, wg *sync.WaitGroup, target *[]config.Skin, list []string, name string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ictx, cancel := chromedp.NewContext(browserCtx)
		defer cancel()

		*target = scrapeList(ictx, list, name)
		fmt.Printf("\033[32m[+]\033[0m Done %s\n", name)
	}()
}

func scrapeList(ctx context.Context, list []string, key string) []config.Skin {
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
				chromedp.Evaluate(string(ConfigJS), nil),
				chromedp.Sleep(config.Delay),
				chromedp.Title(&pageTitle),
				chromedp.Evaluate(string(ScriptJS), &rawData),
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
