package tracker

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go.etcd.io/bbolt"
)

// Define the structure for tracking data
type DayData struct {
	Alcohol   string  `json:"alcohol"`
	Quantity  int     `json:"quantity"`
	Cost      float64 `json:"cost"`
	Timestamp int64   `json:"timestamp"`
}

// Global database instance
var db *bbolt.DB

// Initialize the BoltDB database
func InitDB() error {
	var err error
	db, err = bbolt.Open("tracker.db", 0600, nil) // Creates or opens the database
	if err != nil {
		return err
	}

	fmt.Println("Database initialized")
	return nil
}

// Add a new tracker entry (Hierarchical: Year → Month → Day → Category)
func AddTrackerEntry(year, month, day int, category string, data DayData) error {
	return db.Update(func(tx *bbolt.Tx) error {
		// Create/Get top-level "Tracker" bucket
		root, err := tx.CreateBucketIfNotExists([]byte("Tracker"))
		if err != nil {
			return err
		}

		// Create/Get Year bucket
		yearBucket, err := root.CreateBucketIfNotExists([]byte(fmt.Sprintf("%d", year)))
		if err != nil {
			return err
		}

		// Create/Get Month bucket
		monthBucket, err := yearBucket.CreateBucketIfNotExists([]byte(fmt.Sprintf("%02d", month)))
		if err != nil {
			return err
		}

		// Create/Get Day bucket
		dayBucket, err := monthBucket.CreateBucketIfNotExists([]byte(fmt.Sprintf("%02d", day)))
		if err != nil {
			return err
		}

		// Convert entry to JSON
		entryKey := []byte(category)
		var entries []DayData

		// Check if an entry already exists for this category on this day
		existingData := dayBucket.Get(entryKey)
		if existingData != nil {
			err = json.Unmarshal(existingData, &entries)
			if err != nil {
				return err
			}
		}

		// Append the new data entry
		entries = append(entries, data)

		// Store updated data
		newData, err := json.Marshal(entries)
		if err != nil {
			return err
		}
		fmt.Printf("Entry Added for /%d/%02d/%02d/%s \n ", year, month, day, category)
		return dayBucket.Put(entryKey, newData)
	})
}

// Get all entries for a given year (structured as Year → Month → Day)
func GetEntriesByYear(year int) (map[string]map[string]map[string][]DayData, error) {
	entries := make(map[string]map[string]map[string][]DayData) // Structure: Year -> Month -> Day -> Entries

	err := db.View(func(tx *bbolt.Tx) error {
		root := tx.Bucket([]byte("Tracker"))
		if root == nil {
			return fmt.Errorf("tracker data not found")
		}

		yearBucket := root.Bucket([]byte(fmt.Sprintf("%d", year)))
		if yearBucket == nil {
			return fmt.Errorf("no data found for year %d", year)
		}

		// Iterate over months
		return yearBucket.ForEach(func(monthKey, _ []byte) error {
			monthBucket := yearBucket.Bucket(monthKey)
			if monthBucket == nil {
				return nil
			}

			monthStr := string(monthKey)
			entries[monthStr] = make(map[string]map[string][]DayData)

			// Iterate over days
			return monthBucket.ForEach(func(dayKey, _ []byte) error {
				dayBucket := monthBucket.Bucket(dayKey)
				if dayBucket == nil {
					return nil
				}

				dayStr := string(dayKey)
				entries[monthStr][dayStr] = make(map[string][]DayData)

				// Iterate over categories
				return dayBucket.ForEach(func(categoryKey, value []byte) error {
					var dayEntries []DayData
					if err := json.Unmarshal(value, &dayEntries); err != nil {
						return err
					}

					categoryStr := string(categoryKey)
					entries[monthStr][dayStr][categoryStr] = dayEntries
					return nil
				})
			})
		})
	})

	return entries, err
}

// Get all entries for a given year & month (structured as Month → Day)
func GetEntriesByYearAndMonth(year, month int) (map[string]map[string][]DayData, error) {
	entries := make(map[string]map[string][]DayData) // Structure: Month → Day → Entries

	err := db.View(func(tx *bbolt.Tx) error {
		root := tx.Bucket([]byte("Tracker"))
		if root == nil {
			return fmt.Errorf("tracker data not found")
		}

		yearBucket := root.Bucket([]byte(fmt.Sprintf("%d", year)))
		if yearBucket == nil {
			return fmt.Errorf("no data found for year %d", year)
		}

		monthBucket := yearBucket.Bucket([]byte(fmt.Sprintf("%02d", month)))
		if monthBucket == nil {
			return fmt.Errorf("no data found for month %02d in year %d", month, year)
		}

		// Iterate over days
		return monthBucket.ForEach(func(dayKey, _ []byte) error {
			dayBucket := monthBucket.Bucket(dayKey)
			if dayBucket == nil {
				return nil
			}

			dayStr := string(dayKey)
			entries[dayStr] = make(map[string][]DayData)

			// Iterate over categories
			return dayBucket.ForEach(func(categoryKey, value []byte) error {
				var dayEntries []DayData
				if err := json.Unmarshal(value, &dayEntries); err != nil {
					return err
				}

				categoryStr := string(categoryKey)
				entries[dayStr][categoryStr] = dayEntries
				return nil
			})
		})
	})

	return entries, err
}

// Get all entries for a specific year, month, and day
func GetEntriesByDate(year, month, day int) (map[string][]DayData, error) {
	entries := make(map[string][]DayData) // Structure: Category → Entries

	err := db.View(func(tx *bbolt.Tx) error {
		root := tx.Bucket([]byte("Tracker"))
		if root == nil {
			return fmt.Errorf("tracker data not found")
		}

		yearBucket := root.Bucket([]byte(fmt.Sprintf("%d", year)))
		if yearBucket == nil {
			return fmt.Errorf("no data found for year %d", year)
		}

		monthBucket := yearBucket.Bucket([]byte(fmt.Sprintf("%02d", month)))
		if monthBucket == nil {
			return fmt.Errorf("no data found for month %02d in year %d", month, year)
		}

		dayBucket := monthBucket.Bucket([]byte(fmt.Sprintf("%02d", day)))
		if dayBucket == nil {
			return fmt.Errorf("no data found for day %02d in month %02d of year %d", day, month, year)
		}

		// Iterate over categories
		return dayBucket.ForEach(func(categoryKey, value []byte) error {
			var dayEntries []DayData
			if err := json.Unmarshal(value, &dayEntries); err != nil {
				return err
			}

			categoryStr := string(categoryKey)
			entries[categoryStr] = dayEntries
			return nil
		})
	})

	return entries, err
}

// Get all entries for a specific year, month, and day as a flat list
func GetEntriesByDateList(year, month, day int) ([]DayData, error) {
	allEntries := []DayData{} // Always initialized as an empty slice

	err := db.View(func(tx *bbolt.Tx) error {
		root := tx.Bucket([]byte("Tracker"))
		if root == nil {
			return nil // No data found, return empty list
		}

		yearBucket := root.Bucket([]byte(fmt.Sprintf("%d", year)))
		if yearBucket == nil {
			return nil // No data found, return empty list
		}

		monthBucket := yearBucket.Bucket([]byte(fmt.Sprintf("%02d", month)))
		if monthBucket == nil {
			return nil // No data found, return empty list
		}

		dayBucket := monthBucket.Bucket([]byte(fmt.Sprintf("%02d", day)))
		if dayBucket == nil {
			return nil // No data found, return empty list
		}

		// Iterate over all categories and collect data
		_ = dayBucket.ForEach(func(categoryKey, value []byte) error {
			var entries []DayData
			if err := json.Unmarshal(value, &entries); err != nil {
				fmt.Printf("Warning: Failed to parse category %s: %v\n", string(categoryKey), err)
				return nil // Continue processing other categories even if one fails
			}

			// Append all category entries to the final list
			allEntries = append(allEntries, entries...)
			return nil
		})

		return nil // Always return nil to continue iteration
	})

	// Ensure we return an empty slice, never nil
	if err != nil {
		return []DayData{}, err
	}

	return allEntries, nil
}

// Get all entries for a specific year, month, day, and category
func GetEntriesByDateCategory(year, month, day int, category string) ([]DayData, error) {
	var entries []DayData

	err := db.View(func(tx *bbolt.Tx) error {
		root := tx.Bucket([]byte("Tracker"))
		if root == nil {
			return fmt.Errorf("tracker data not found")
		}

		yearBucket := root.Bucket([]byte(fmt.Sprintf("%d", year)))
		if yearBucket == nil {
			return fmt.Errorf("no data found for year %d", year)
		}

		monthBucket := yearBucket.Bucket([]byte(fmt.Sprintf("%02d", month)))
		if monthBucket == nil {
			return fmt.Errorf("no data found for month %02d in year %d", month, year)
		}

		dayBucket := monthBucket.Bucket([]byte(fmt.Sprintf("%02d", day)))
		if dayBucket == nil {
			return fmt.Errorf("no data found for day %02d in month %02d of year %d", day, month, year)
		}

		// Retrieve the specific category data
		value := dayBucket.Get([]byte(category))
		if value == nil {
			return fmt.Errorf("no data found for category '%s' on %02d-%02d-%d", category, day, month, year)
		}

		// Unmarshal JSON into entries slice
		if err := json.Unmarshal(value, &entries); err != nil {
			return err
		}

		return nil
	})

	return entries, err
}

// PrintDaysInYear prints all the days present in the database for a given year
// PrintDaysWithDataInYear prints all days with their associated DayData entries for a given year
func PrintDaysInYear(year int) error {
	return db.View(func(tx *bbolt.Tx) error {
		root := tx.Bucket([]byte("Tracker"))
		if root == nil {
			fmt.Println("No tracker data found.")
			return nil
		}

		yearBucket := root.Bucket([]byte(fmt.Sprintf("%d", year)))
		if yearBucket == nil {
			fmt.Printf("No data found for year %d\n", year)
			return nil
		}

		// Iterate over months
		return yearBucket.ForEach(func(monthKey, _ []byte) error {
			monthBucket := yearBucket.Bucket(monthKey)
			if monthBucket == nil {
				return nil
			}

			// Iterate over days
			return monthBucket.ForEach(func(dayKey, _ []byte) error {
				dayBucket := monthBucket.Bucket(dayKey)
				if dayBucket == nil {
					return nil
				}

				fmt.Printf("Entries for date: %d-%s-%s\n", year, string(monthKey), string(dayKey))

				// Iterate over categories and print DayData entries
				return dayBucket.ForEach(func(categoryKey, value []byte) error {
					var entries []DayData
					if err := json.Unmarshal(value, &entries); err != nil {
						fmt.Printf("Error parsing data for category %s: %v\n", string(categoryKey), err)
						return nil // Continue processing other categories even if one fails
					}

					fmt.Printf("  Category: %s\n", string(categoryKey))
					for _, entry := range entries {
						fmt.Printf("    Alcohol: %s, Quantity: %d, Cost: %.2f, Timestamp: %d\n",
							entry.Alcohol, entry.Quantity, entry.Cost, entry.Timestamp)
					}
					return nil
				})
			})
		})
	})
}

// FindLatestEntryDate retrieves the latest (most recent) entry from the database.
func FindLatestEntryDate() (int, int, int, error) {
	var latestYear, latestMonth, latestDay int

	err := db.View(func(tx *bbolt.Tx) error {
		root := tx.Bucket([]byte("Tracker"))
		if root == nil {
			return fmt.Errorf("tracker data not found")
		}

		// Iterate over years in reverse order to find the most recent year
		err := root.ForEach(func(yearKey, _ []byte) error {
			year, _ := strconv.Atoi(string(yearKey))
			if year > latestYear {
				latestYear = year
			}
			return nil
		})
		if err != nil {
			return err
		}
		if latestYear == 0 {
			return fmt.Errorf("no data found in database")
		}

		yearBucket := root.Bucket([]byte(fmt.Sprintf("%d", latestYear)))
		if yearBucket == nil {
			return fmt.Errorf("no data found for latest year %d", latestYear)
		}

		// Iterate over months in reverse order to find the most recent month
		err = yearBucket.ForEach(func(monthKey, _ []byte) error {
			month, _ := strconv.Atoi(string(monthKey))
			if month > latestMonth {
				latestMonth = month
			}
			return nil
		})
		if err != nil {
			return err
		}

		monthBucket := yearBucket.Bucket([]byte(fmt.Sprintf("%02d", latestMonth)))
		if monthBucket == nil {
			return fmt.Errorf("no data found for latest month %02d in year %d", latestMonth, latestYear)
		}

		// Iterate over days in reverse order to find the most recent day
		err = monthBucket.ForEach(func(dayKey, _ []byte) error {
			day, _ := strconv.Atoi(string(dayKey))
			if day > latestDay {
				latestDay = day
			}
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, 0, 0, err
	}

	return latestYear, latestMonth, latestDay, nil
}

// GetDaysSinceLastEntry calculates the number of days since the latest entry.
func GetDaysSinceLastEntry() (int, error) {
	latestYear, latestMonth, latestDay, err := FindLatestEntryDate()
	if err != nil {
		return -1, err
	}

	latestDate := time.Date(latestYear, time.Month(latestMonth), latestDay, 0, 0, 0, 0, time.UTC)
	today := time.Now().UTC()
	daysSince := int(today.Sub(latestDate).Hours() / 24)

	return daysSince, nil
}

// DeleteEntry removes an entry based on the given year, month, day, category, and timestamp.
func DeleteEntry(year, month, day int, category string, timestamp int64) error {
	return db.Update(func(tx *bbolt.Tx) error {
		root := tx.Bucket([]byte("Tracker"))
		if root == nil {
			return fmt.Errorf("tracker data not found")
		}

		yearBucket := root.Bucket([]byte(fmt.Sprintf("%d", year)))
		if yearBucket == nil {
			return fmt.Errorf("no data found for year %d", year)
		}

		monthBucket := yearBucket.Bucket([]byte(fmt.Sprintf("%02d", month)))
		if monthBucket == nil {
			return fmt.Errorf("no data found for month %02d in year %d", month, year)
		}

		dayBucket := monthBucket.Bucket([]byte(fmt.Sprintf("%02d", day)))
		if dayBucket == nil {
			return fmt.Errorf("no data found for day %02d in month %02d of year %d", day, month, year)
		}

		// Retrieve the specific category data
		entryKey := []byte(category)
		value := dayBucket.Get(entryKey)
		if value == nil {
			return fmt.Errorf("no data found for category '%s' on %02d-%02d-%d", category, day, month, year)
		}

		// Unmarshal JSON into entries slice
		var entries []DayData
		if err := json.Unmarshal(value, &entries); err != nil {
			return err
		}

		// Filter out the entry with the matching timestamp
		filteredEntries := []DayData{}
		for _, entry := range entries {
			if entry.Timestamp != timestamp {
				filteredEntries = append(filteredEntries, entry)
			}
		}

		// If no entries remain, delete the category key
		if len(filteredEntries) == 0 {
			return dayBucket.Delete(entryKey)
		}

		// Store updated data
		newData, err := json.Marshal(filteredEntries)
		if err != nil {
			return err
		}

		return dayBucket.Put(entryKey, newData)
	})
}

// Close the database connection
func CloseDB() {
	if db != nil {
		db.Close()
	}
}
