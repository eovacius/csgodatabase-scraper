// main configuration file where you can set options for scraper,
// add/remove subdomains, change constants for scraper settings etc.

package config

import (
	"time"

	"github.com/chromedp/chromedp"
)

// here you should add subdomains (in our case it's name of collection/case)
var List = []string{
	"kilowatt-case", "revolution-case", "recoil-case", "dreams-nightmares-case", "sealed-genesis-terminal",
	"snakebite-case", "fracture-case", "prisma-2-case", "cs20-case", "prisma-case", "danger-zone-case",
	"horizon-case", "clutch-case", "spectrum-2-case", "operation-hydra-case", "spectrum-case", "glove-case",
	"gamma-2-case", "gamma-case", "chroma-3-case", "operation-wildfire-case", "revolver-case", "shadow-case", "falchion-case",
	"chroma-2-case", "chroma-case", "operation-vanguard-weapon-case", "operation-breakout-weapon-case", "huntsman-weapon-case",
	"operation-phoenix-weapon-case", "csgo-weapon-case-3", "winter-offensive-weapon-case", "csgo-weapon-case-2", "operation-bravo-case",
	"csgo-weapon-case", "fever-case", "gallery-case", "operation-riptide-case", "operation-broken-fang-case", "shattered-web-case", "esports-2014-summer-case",
	"esports-2013-winter-case", "esports-2013-case", "anubis-collection-package", "x-ray-p250-package",
}

// collections to scrape
var CollectionsList = []string{
	"the-2018-inferno-collection", "the-2018-nuke-collection", "the-ascent-collection", "the-boreal-collection", "the-dust-2-collection",
	"the-radiant-collection", "the-safehouse-collection", "limited-edition-item", "the-graphic-design-collection", "the-overpass-2024-collection",
	"the-sport-field-collection", "the-train-2025-collection", "the-2021-dust-2-collection", "the-2021-mirage-collection", "the-2021-train-collection",
	"the-2021-vertigo-collection", "the-ancient-collection", "the-control-collection", "the-havoc-collection", "the-canals-collection", "the-norse-collection",
	"the-st-marc-collection", "the-cache-collection", "the-chop-shop-collection", "the-cobblestone-collection", "the-gods-and-monsters-collection", "the-overpass-collection",
	"the-rising-sun-collection", "the-anubis-collection", "the-assault-collection", "the-aztec-collection", "the-baggage-collection", "the-bank-collection",
	"the-dust-collection", "the-inferno-collection", "the-italy-collection", "the-lake-collection", "the-militia-collection", "the-mirage-collection", "the-nuke-collection",
	"the-office-collection", "the-train-collection", "the-vertigo-collection", "the-alpha-collection", "the-blacksite-collection",
}

// Scraper settings
var (
	Target         = "https://www.csgodatabase.com" // target site to scrape from
	DeadLine       = 120 * time.Second              // time limit for context
	UrlLengthLimit = 60                             // shorten url to specified length
	Delay          = 1 * time.Second                // delay to avoid triggering site protections
	Headless       = true                           // run browser in headless/headed mode
)

// allocator options
var Opts = append(chromedp.DefaultExecAllocatorOptions[:],
	chromedp.Flag("headless", Headless),
	chromedp.Flag("disable-blink-features", "AutomationControlled"),
	chromedp.Flag("blink-settings", "imagesEnabled=false"),
	chromedp.Flag("exclude-switches", "enable-automation"),
	chromedp.Flag("disable-extensions", false),
	chromedp.Flag("start-maximized", false),
	chromedp.Flag("window-size", "800,600"),
	chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36"),
)
