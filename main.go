package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/eovacius/csgodatabase-scraper/scraper"
	"github.com/eovacius/csgodatabase-scraper/scraper/config"
)

func main() {
	fmt.Println("[*] Starting scraper...")

	skins, err := scraper.ScrapeSkins()
	if err != nil {
		log.Fatalf("[!] Error during scraping: %v", err)
	}

	if len(skins) == 0 {
		log.Fatalf("[!] No skins were scraped! Exiting...")
	}

	jsonData, err := json.MarshalIndent(skins, "", "  ")
	if err != nil {
		log.Fatalf("[!] Error marshaling JSON: %v", err)
	}

	// if you want data to be printed to console
	if config.ConsoleLog {
		fmt.Println(string(jsonData))
	}
	// generate filename with current date
	timestamp := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("json/data_%s.json", timestamp)

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		log.Fatalf("[!] Error writing to file: %v", err)
	}

	fmt.Println("\n[+] Done. See files inside json folder")
}
