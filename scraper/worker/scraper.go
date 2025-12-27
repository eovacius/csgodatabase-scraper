package worker

import (
	"context"
	_ "embed"
	"fmt"
	"sync"

	"github.com/chromedp/chromedp"
	"github.com/eovacius/csgodatabase-scraper/internal"
	"github.com/eovacius/csgodatabase-scraper/scraper/agentscraper"
	"github.com/eovacius/csgodatabase-scraper/scraper/config"
	skinsscraper "github.com/eovacius/csgodatabase-scraper/scraper/skinscraper"
)

// Run func executes the scraping process and returns a slice of Skin structs with error if any.
func ScrapeSkins() ([]config.Skin, []config.Agent, error) {
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
	var agents []config.Agent

	spawnWorker(browserCtx, &wg, &cases, config.List, "cases")
	spawnWorker(browserCtx, &wg, &collections, config.CollectionsList, "collections")
	spawnAgentWorker(browserCtx, &wg, &agents, config.Agents, "agents")

	// thats to test whether we get detected or not
	// spawnWorker(browserCtx, &wg, &cases, config.List, "cases")
	// spawnWorker(browserCtx, &wg, &collections, config.CollectionsList, "collections")

	wg.Wait()

	merged := append(cases, collections...)

	uniqueSkins := internal.RemoveDuplicates(merged)
	return uniqueSkins, agents, nil
}

func spawnAgentWorker(browserCtx context.Context, wg *sync.WaitGroup, target *[]config.Agent, list []string, name string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ictx, cancel := chromedp.NewContext(browserCtx)
		defer cancel()

		*target = agentscraper.ScrapeAgentsList(ictx, list, name)
		fmt.Printf("\033[32m[+]\033[0m Done %s\n", name)
	}()
}

func spawnWorker(browserCtx context.Context, wg *sync.WaitGroup, target *[]config.Skin, list []string, name string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ictx, cancel := chromedp.NewContext(browserCtx)
		defer cancel()

		*target = skinsscraper.ScrapeList(ictx, list, name)
		fmt.Printf("\033[32m[+]\033[0m Done %s\n", name)
	}()
}
