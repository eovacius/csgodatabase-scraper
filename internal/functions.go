package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/eovacius/csgodatabase-scraper/scraper/config"
)

func RemoveDuplicates(skins []config.Skin) []config.Skin {
	seen := make(map[string]bool)
	var unique []config.Skin

	for _, s := range skins {
		key := s.Name + "|" + s.Weapon + "|" + s.Rarity
		if !seen[key] {
			seen[key] = true
			unique = append(unique, s)
		}
	}

	// temporary filter to remove souvenir packages as scraper can't handle them without separating each souvenir by subdomain
	var filtered []config.Skin
	for _, skin := range unique {
		if skin.Weapon != "Souvenir Package" {
			filtered = append(filtered, skin)
		}
	}
	return filtered
}

func ParsePrice(raw, stattRaw string) config.Price {
	raw = strings.TrimSpace(raw)
	stattrakRaw := strings.TrimSpace(stattRaw)

	var currency string
	if strings.Contains(raw, "$") || strings.Contains(stattrakRaw, "$") {
		currency = "USD"
	}

	price := config.Price{
		PriceString:         raw,
		PriceStattrakString: stattrakRaw,
		Currency:            currency,
		Min:                 config.PriceValue{Value: 0, StattrakValue: 0, Unit: currency},
		Max:                 config.PriceValue{Value: 0, StattrakValue: 0, Unit: currency},
		UpdatedAt:           time.Now().Format(time.RFC3339),
	}

	re := regexp.MustCompile(`[\d.,]+`)
	if raw != "" {
		parts := strings.Split(raw, "-")
		if len(parts) == 1 {
			v, _ := strconv.ParseFloat(strings.ReplaceAll(re.FindString(parts[0]), ",", ""), 64)
			price.Min.Value = v
			price.Max.Value = v
		} else if len(parts) >= 2 {
			v1, _ := strconv.ParseFloat(strings.ReplaceAll(re.FindString(parts[0]), ",", ""), 64)
			v2, _ := strconv.ParseFloat(strings.ReplaceAll(re.FindString(parts[1]), ",", ""), 64)
			price.Min.Value = v1
			price.Max.Value = v2
		}
	}

	if stattrakRaw != "" {
		parts := strings.Split(stattrakRaw, "-")
		if len(parts) == 1 {
			v, _ := strconv.ParseFloat(strings.ReplaceAll(re.FindString(parts[0]), ",", ""), 64)
			price.Min.StattrakValue = v
			price.Max.StattrakValue = v
		} else if len(parts) >= 2 {
			v1, _ := strconv.ParseFloat(strings.ReplaceAll(re.FindString(parts[0]), ",", ""), 64)
			v2, _ := strconv.ParseFloat(strings.ReplaceAll(re.FindString(parts[1]), ",", ""), 64)
			price.Min.StattrakValue = v1
			price.Max.StattrakValue = v2
		}
	}

	return price
}

func SpecialMark(weapon string) string {
	keywords := []string{"Knife", "Gloves", "Wraps"}
	for _, keyword := range keywords {
		if strings.Contains(weapon, keyword) {
			return "★ " + weapon
		}
	}
	return weapon
}

func SaveJSON(path string, data interface{}) {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("[!] Failed to marshal JSON: %v\n", err)
		return
	}

	err = os.WriteFile(path, file, 0644)
	if err != nil {
		fmt.Printf("[!] Failed to write file %s: %v\n", path, err)
		return
	}

	fmt.Printf("[+] Saved: %s\n", path)
}
