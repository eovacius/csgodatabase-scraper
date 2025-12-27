package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/eovacius/csgodatabase-scraper/scraper/config"
	"github.com/eovacius/csgodatabase-scraper/scraper/worker"
)

func main() {
	// cli flags
	aggressive := flag.Bool("aggressive", false, "Run scraper aggressively (less delay)")
	stealth := flag.Bool("stealth", false, "Run scraper in stealth mode (randomized and more human like delay)")

	flag.Parse()

	if *aggressive {
		config.Delay = 0 * time.Millisecond
		fmt.Println("Aggressive mode. Delay:", config.Delay)
	}

	if *stealth {
		randomMs := 500 + rand.Intn(1500)
		config.Delay += time.Duration(randomMs) * time.Millisecond
	}

	fmt.Println("[*] Starting scraper...")

	skins, agents, err := worker.ScrapeSkins()
	if err != nil {
		log.Fatalf("\033[31m[!]\033[0m Error during scraping: %v", err)
	}

	combinedData := config.DataOutput{
		Skins:  skins,
		Agents: agents,
	}

	if len(skins) == 0 {
		log.Fatalf("\033[31m[!]\033[0m No skins were scraped! Exiting...")
	}

	jsonData, err := json.MarshalIndent(combinedData, "", "  ")
	if err != nil {
		log.Fatalf("\033[31m[!]\033[0m Error marshaling JSON: %v", err)
	}

	// generate filename with current date
	timestamp := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("json/data_%s.json", timestamp)

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		log.Fatalf("\033[31m[!]\033[0m Error writing to file: %v", err)
	}

	fmt.Println("\n\033[32m[+]\033[0m Done. See files inside json folder")
}
