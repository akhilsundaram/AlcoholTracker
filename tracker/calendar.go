package tracker

import (
	"errors"
	"fmt"
	"time"
)

func ValidateDate(day, month, year int) error {
	if year < 2000 || year > 2100 {
		return errors.New("year must be between 2000 and 2100")
	}
	date := fmt.Sprintf("%d-%02d-%02d", year, month, day)
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("invalid date: %v", err)
	}
	return nil
}
