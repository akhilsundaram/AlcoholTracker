package main

import (
	"AlcoholTracker/tracker"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	DBerr := tracker.InitDB()
	if DBerr != nil {
		log.Fatal("Failed to initialize database:", DBerr)
	}
}

// Shutdown
func (a *App) Shutdown(ctx context.Context) {
	tracker.CloseDB()
	log.Println("Database closed.")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Expose AddTrackerEntry to the frontend
func (a *App) AddTrackerEntry(year int, month int, day int, category string, quantity int, cost float64) {
	entry := tracker.DayData{
		Alcohol:   category,
		Quantity:  quantity,
		Cost:      cost,
		Timestamp: time.Now().Unix(),
	}

	err := tracker.AddTrackerEntry(year, month, day, category, entry)
	if err != nil {
		runtime.LogError(a.ctx, "Error adding entry: "+err.Error())
	}
}

// Expose GetEntriesByDate to the frontend
func (a *App) GetEntriesByDate(year int, month int, day int, category string) []tracker.DayData {
	entries, err := tracker.GetEntriesByDateList(year, month, day)
	if err != nil {
		runtime.LogError(a.ctx, "Error fetching entries: "+err.Error())
	}
	return entries
}

func (a *App) GetAlcoholCategories() []string {
	return tracker.GetAlcoholTypes()
}

func (a *App) ValidateFormDate(year, month, day int) bool {
	err := tracker.ValidateDate(day, month, year)
	return err == nil
}

func (a *App) GetDrinks(year, month, day int) string {
	categories := []string{"empty", "low", "moderate", "heavy", "binge", "excessive"}
	drinks, err := tracker.GetTotalDrinksOnDay(year, month, day)
	if err != nil {
		fmt.Println(drinks, err)
		return categories[0]
	}

	switch {
	case drinks == -1:
		return categories[0]
	case drinks >= 0 && drinks < 1:
		return categories[1]
	case drinks >= 1 && drinks < 2:
		return categories[2]
	case drinks >= 2 && drinks < 3:
		return categories[3]
	case drinks >= 3 && drinks < 5:
		return categories[4]
	case drinks >= 5:
		return categories[5]
	default:
		return categories[0]
	}
}

func (a *App) GetDrinkTagColor(year, month, day int) int {
	// categories := []string{"gray", "#60aa9b", "#43766c", "#ffdf60", "#fa8072", "#ed4d09"}
	drinks, err := tracker.GetTotalDrinksOnDay(year, month, day)
	if err != nil {
		fmt.Println(drinks, err)
		return 0
	}

	switch {
	case drinks == -1:
		return 0
	case drinks >= 0 && drinks < 1:
		return 1
	case drinks >= 1 && drinks < 2:
		return 2
	case drinks >= 2 && drinks < 3:
		return 3
	case drinks >= 3 && drinks < 5:
		return 4
	case drinks >= 5:
		return 5
	default:
		return 0
	}
}

func (a *App) GetDrinkCount(year, month, day int) float64 {
	drinks, err := tracker.GetTotalDrinksOnDay(year, month, day)
	if err != nil {
		fmt.Println(drinks, err)
		return 0
	}

	switch {
	case drinks == -1:
		return 0
	default:
		return drinks
	}
}

func (a *App) GetEntriesOnDate(year, month, day int) (map[string][]tracker.DayData, error) {
	entries, err := tracker.GetEntriesByDate(year, month, day)
	if err != nil {
		fmt.Println(err)
	}
	return entries, err
}

func (a *App) GetDaysSinceLastDrink() int {
	days, err := tracker.GetDaysSinceLastEntry()
	if err != nil {
		fmt.Println(days, err)
	}
	fmt.Println(days)
	return days
}

func (a *App) DeleteDrink(year, month, day int, alcohol string, timestamp int64) bool {
	err := tracker.DeleteEntry(year, month, day, alcohol, timestamp)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

// Expose AddTrackerEntry to the frontend
func (a *App) AddTrackerEntryUpdate(year int, month int, day int, category string, quantity int, cost float64, timestamp int64) {
	entry := tracker.DayData{
		Alcohol:   category,
		Quantity:  quantity,
		Cost:      cost,
		Timestamp: timestamp,
	}

	err := tracker.AddTrackerEntry(year, month, day, category, entry)
	if err != nil {
		runtime.LogError(a.ctx, "Error adding entry: "+err.Error())
	}
}
