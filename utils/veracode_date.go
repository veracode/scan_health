package utils

import (
	"fmt"
	"strings"
	"time"
)

func ParseVeracodeDate(date string) time.Time {
	var format = "2006-01-02 15:04:05 MST"

	if (!strings.Contains(date, "UTC")) && strings.Contains(date, "T") {
		format = "2006-01-02T15:04:05-07:00"
	}

	parsed, err := time.Parse(format, date)

	if err != nil {
		ErrorAndExit(fmt.Sprintf("Could not parse \"%s\" as a date", date), err)
	}

	return parsed
}
