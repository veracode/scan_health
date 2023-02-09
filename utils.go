package main

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func colorPrintf(format string) {
	color.New().Printf(format)
}

func isStringInStringArray(input string, list []string) bool {
	for _, item := range list {
		if input == item {
			return true
		}
	}

	return false
}

func stringToFloat(input string) (float64, error) {
	return strconv.ParseFloat(input, 64)
}

func getSortedIntArrayAsFormattedString(list []int) string {
	sort.Ints(list[:])
	var output []string
	for _, x := range list {
		output = append(output, strconv.Itoa(x))
	}

	return strings.Join(output, ",")
}

func isInIntArray(x int, y []int) bool {
	for _, z := range y {
		if x == z {
			return true
		}
	}

	return false
}

func getFormattedOnlyInSideString(side string) string {
	if side == "A" {
		return color.HiGreenString("Only in A")
	}

	return color.HiMagentaString("Only in B")
}

func getFormattedSideString(side string) string {
	if side == "A" {
		return color.HiGreenString("A")
	}

	return color.HiMagentaString("B")
}

func getFormattedSideStringWithMessage(side, message string) string {
	if side == "A" {
		return color.HiGreenString(message)
	}

	return color.HiMagentaString(message)
}

func dedupeArray[T interface{}](array []T) []T {
	result := []T{}

	for _, item := range array {
		found := false
		for _, processedItem := range result {
			if !found && reflect.DeepEqual(item, processedItem) {
				found = true
			}
		}

		if !found {
			result = append(result, item)
		}
	}

	return result
}

func formatDuration(duration time.Duration) string {
	str := duration.String()

	if duration.Hours() > 24 {
		days := int(duration.Hours() / 24)
		remainingHours := duration.Hours() - float64(days*24)
		str = fmt.Sprintf("%dd %dh%s", days, int(remainingHours), strings.Split(str, "h")[1])
	}

	str = strings.Replace(str, "h", "h ", 1)
	str = strings.Replace(str, "m", "m ", 1)

	return str
}
