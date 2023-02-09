package main

import (
	"os"
	"time"

	"github.com/fatih/color"
)

func parseVeracodeDate(date string) time.Time {
	parsed, err := time.Parse("2006-01-02 15:04:05 MST", date)

	if err != nil {
		color.HiRed("Error: Could not parse \"%s\" as a date", date)
		os.Exit(1)
	}

	return parsed
}
