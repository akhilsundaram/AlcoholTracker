package tracker

import (
	"time"

	"github.com/dgraph-io/badger"
)

// Alcohol types with their Alcohol By Volume (ABV) percentages
var alcoholMap = map[string]float64{
	"Beer":    5.0,  // 5% ABV
	"Soju":    16.0, // 16% ABV
	"Wine":    12.0, // 12% ABV
	"Vodka":   40.0, // 40% ABV
	"Rum":     40.0, // 40% ABV
	"Whiskey": 40.0, // 40% ABV
	"Gin":     37.5, // 37.5% ABV
}

// getAlcoholTypes returns a slice of all the keys in alcoholMap
func GetAlcoholTypes() []string {
	keys := make([]string, 0, len(alcoholMap))
	for key := range alcoholMap {
		keys = append(keys, key)
	}
	return keys
}

// Function to calculate standard drinks
func CalculateStandardDrinks(volumeML float64, alcoholType string) float64 {
	abv, exists := alcoholMap[alcoholType]
	standardDrinks := 0.0
	// Formula: (volume in mL * ABV%) / 17.7
	if !exists {
		standardDrinks = (volumeML * (40.0 / 100)) / 17.7
	} else {
		standardDrinks = (volumeML * (abv / 100)) / 17.7
	}
	return standardDrinks
}

// GetTotalDrinksOnDay calculates the total number of standard drinks consumed on a given day
func GetTotalDrinksOnDay(year, month, day int) (float64, error) {
	totalDrinks := 0.0
	drinksFound := false

	for category := range alcoholMap {
		entries, err := GetEntriesByDateCategory(year, month, day, category)
		if err != nil {
			// If key doesn't exist, skip this category
			if err == badger.ErrKeyNotFound {
				continue
			}

		}

		if len(entries) > 0 {
			drinksFound = true
		}

		for _, entry := range entries {
			// Assuming quantity represents volume in mL
			totalDrinks += CalculateStandardDrinks(float64(entry.Quantity), category)
		}
	}

	// If no drinks were found, return -1 and nil error
	if !drinksFound {
		return -1, nil
	}

	return totalDrinks, nil
}

func GetTotalDrinksToday() (float64, error) {
	// Get current date
	now := time.Now()
	year, month, day := now.Year(), int(now.Month()), now.Day()

	// Get total drinks for today
	totalDrinks, err := GetTotalDrinksOnDay(year, month, day)
	if err != nil {
		return -1, err
	}

	if totalDrinks == -1 {
		return 0, nil
	}

	return totalDrinks, nil
}
