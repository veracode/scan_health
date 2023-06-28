package utils

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func ColorPrintf(format string) {
	_, _ = color.New().Print(format)
}

func IsStringInStringArray(input string, list []string) bool {
	for _, item := range list {
		if input == item {
			return true
		}
	}

	return false
}

func StringToFloat(input string) (float64, error) {
	return strconv.ParseFloat(input, 64)
}

func DedupeArray[T interface{}](array []T) []T {
	var result []T

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

func FormatDuration(duration time.Duration) string {
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

func PrintTitle(title string) {
	color.HiCyan(title)
	fmt.Println(strings.Repeat("=", len(title)))
}

func Top5StringList(items []string) string {
	sort.Strings(items)

	if len(items) > 5 {
		return fmt.Sprintf("\"%s\" and %d others", strings.Join(items[0:5], "\", \""), len(items)-5)
	}

	return fmt.Sprintf("\"%s\"", strings.Join(items, "\", \""))
}

func ErrorAndExit(message string, err error) {
	color.HiRed(fmt.Sprintf("Error: %s\n", message))

	if err != nil {
		fmt.Printf("%v\n", err)
	}

	os.Exit(1)
}

func ErrorAndExitWithUsage(message string) {
	color.HiRed(fmt.Sprintf("Error: %s", message))
	print("\nUsage:\n")
	flag.PrintDefaults()
	os.Exit(1)
}
