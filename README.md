## Cs2Scraper  [![PkgGoDev](https://pkg.go.dev/badge/github.com/eovacius/csgodatabase-scraper)](https://pkg.go.dev/github.com/eovacius/csgodatabase-scraper) (Beta)


<strong>csgodatabase-scraper</strong> is a scraping tool to retrieve fresh data with prices and wrap it up in convenient JSON file with clean structure
Example json:

```json
{
    "name": "Inheritance",
    "weapon": "AK-47",
    "rarity": "Covert Rifle",
    "collection": "The Kilowatt Collection",
    "price": {
      "price_string": "$48.62 - $302.37",
      "price_stattrak_string": "$91.61 - $694.48",
      "currency": "USD",
      "min": {
        "value": 48.62,
        "stattrak_value": 91.61,
        "unit": "USD"
      },
      "max": {
        "value": 302.37,
        "stattrak_value": 694.48,
        "unit": "USD"
      },
      "updated_at": "2025-11-09T20:31:26+04:00"
    },
    "url": "https://www.example.com/images/AK-47_Inheritance.png"
}
```

## Script offers:

1. **Structured data** Pre‑scraped and well‑organized JSON files available in the `json` folder.
2. **Fresh scraping** Run the script anytime to fetch up‑to‑date skin data and latest prices.
3. **Configuration** Flexible options stored in the `config` folder for easy adjustments.
4. **Speed** Uses goroutines to scrape asynchronously and significantly reduce total runtime.
5. **Stealth** Built‑in safe delays ensure Cloudflare is not triggered under default settings.
6. **JS proof** 100% reliable scraping even on JS heavy site like <strong>csgo-database</strong>.


<strong>Disclaimer!</strong>

Since the script subtly “unofficially” targets online databases, in our case the popular site <a href="https://www.csgodatabase.com/">CSGO Database</a>, you must understand that all responsibility lies with you. Whether you use the data for a hobby project, a real website, or any form of monetization, you are solely accountable for the consequences.

Even though the script simply retrieves data that is already visible to you, but script is developed the way to bypass anti-bot systems, which can be considered a violation of the target website’s terms of service.

## About Technical Stuff
The scraper now works **asynchronously** using goroutines (not as stealth but still) by parsing data from different tabs with isolated contextes within one single chromedp instance for memory efficiency. 

The script uses a Chromium instance to make humanized requests to pre-configured URLs, injects custom JavaScript to retrieve data specified by selectors, and then navigates to the next URL in sequence. This approach ensures reliable scraping while minimizing detection.

## Why Chromium?
The scraper uses the `chromedp` library because it provides a **high-level interface for controlling Chromium** programmatically. This allows us to:

- Navigate websites like a real browser.
- Execute custom JavaScript on the page to extract dynamic content.
- Bypass certain basic anti-bot protections that block simple HTTP requests.
- Wait for elements to load and handle asynchronous content reliably.
- And it helps disguise the script as a real user because it operates through an actual browser.

## What you'll need

1. <a href="https://go.dev/doc/install">Golang</a> obviously
2. Chromium `sudo apt install chromium-browser`


## Usage

1. **Clone repository**

```bash
git clone https://github.com/eovacius/cs2scraper.git
```

2. **Ensure that dependencies are installed**

```bash
go mod tidy
```

3. **Run script in root**

```bash
go run .
```















