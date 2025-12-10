package scraper

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"sync"

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
	fmt.Print(`
   ____  ____     ____    ____      ____    ____        _       ____   U _____ u   ____     
U /"___|/ __"| u |___"\  / __"| uU /"___|U |  _"\ u U  /"\  u U|  _"\ u\| ___"|/U |  _"\ u  
\| | u <\___ \/  U __) |<\___ \/ \| | u   \| |_) |/  \/ _ \/  \| |_) |/ |  _|"   \| |_) |/  
 | |/__ u___) |  \/ __/ \u___) |  | |/__   |  _ <    / ___ \   |  __/   | |___    |  _ <    
  \____||____/>> |_____|u|____/>>  \____|  |_| \_\  /_/   \_\  |_|      |_____|   |_| \_\   
 _// \\  )(  (__)<<  //   )(  (__)_// \\   //   \\_  \\    >>  ||>>_    <<   >>   //   \\_  
(__)(__)(__)    (__)(__) (__)    (__)(__) (__)  (__)(__)  (__)(__)__)  (__) (__) (__)  (__) 
		`)

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
		fmt.Printf("[+] done %s\n", name)
	}()
}

func scrapeList(ctx context.Context, list []string, key string) []config.Skin {
	var allSkins []config.Skin
	path := "json/" + key
	_ = os.Mkdir(path, 0755)

	for _, item := range list {
		url := config.Target + "/" + key + "/" + item + "/"

		fmt.Printf("Trying to scrape: %s\n", url)

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
			fmt.Printf("\nchromedp error: %v\n", err)
			continue
		}

		var skins []config.Skin
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
			fmt.Printf("[!] No skins found for %s\n", item)
			continue
		}

		filename := fmt.Sprintf("%s/%s.json", path, item)
		internal.SaveJSON(filename, skins)

		allSkins = append(allSkins, skins...)
	}

	return allSkins
}
