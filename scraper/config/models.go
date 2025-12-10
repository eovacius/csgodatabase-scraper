package config

// all those structs are main models of our json structure, example in structure.json
type Skin struct {
	Name       string `json:"name"`
	Weapon     string `json:"weapon"`
	Rarity     string `json:"rarity"`
	Collection string `json:"collection"`
	Price      Price  `json:"price"`
	URL        string `json:"url"`
}

type Price struct {
	PriceString         string     `json:"price_string"`
	PriceStattrakString string     `json:"price_stattrak_string"`
	Currency            string     `json:"currency"`
	Min                 PriceValue `json:"min"`
	Max                 PriceValue `json:"max"`
	UpdatedAt           string     `json:"updated_at"`
}

type PriceValue struct {
	Value         float64 `json:"value"`
	StattrakValue float64 `json:"stattrak_value"`
	Unit          string  `json:"unit"`
}
