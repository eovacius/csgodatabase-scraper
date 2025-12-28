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
	// fancy ascii title ;) happy 2026!
	fmt.Print("\033[34m" + `
                                .___       __        ___.                                                                                  
  ____   ______ ____   ____   __| _/____ _/  |______ \_ |__ _____    ______ ____             ______ ________________  ______   ___________ 
_/ ___\ /  ___// ___\ /  _ \ / __ |\__  \\   __\__  \ | __ \\__  \  /  ___// __ \   ______  /  ___// ___\_  __ \__  \ \____ \_/ __ \_  __ \
\  \___ \___ \/ /_/  >  <_> ) /_/ | / __ \|  |  / __ \| \_\ \/ __ \_\___ \\  ___/  /_____/  \___ \\  \___|  | \// __ \|  |_> >  ___/|  | \/  
 \___  >____  >___  / \____/\____ |(____  /__| (____  /___  (____  /____  >\___  >         /____  >\___  >__|  (____  /   __/ \___  >__|   
     \/     \/_____/             \/     \/          \/    \/     \/     \/     \/               \/     \/           \/|__|        \/     
		
✻     °      ❅      ⁎     ❄        +  ❆       ° ❉  °  ✽  ⁎     ❉  +    °      ✻  ❆            ❅❄        ✽      °   ✼  ❄         ⁎ 
❅ °    ✼  ✻     ⁎       ✽       +      ❆ ✻     °      ❅      ⁎     ❄        +  ❆       ° ❉°  ✽  ⁎     ❉  +    °      ✻  ❆            ❅
    ❄        ✽      °   ✼  ❄         ⁎     ❅°    ✼  ✻     ⁎       ✽       +      ❆  ° ❉°  ✽  ⁎     ❉  +    °      ✻  ❆            ❅
	 
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
	uniqueAgents := internal.RemoveAgentDuplicates(agents)

	return uniqueSkins, uniqueAgents, nil

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
