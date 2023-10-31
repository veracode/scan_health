package checks

import (
	"fmt"
	"regexp"

	"github.com/antfie/scan_health/v2/report"
)

func UNUSED(x ...interface{}) {}

func compareModuleCount(r *report.Report, pr *report.Report) {

	previousModuleCount := len(pr.Modules)
	currentModuleCount := len(r.Modules)

	if previousModuleCount != currentModuleCount {

		previousModuleWord := "module"
		if previousModuleCount != 1 {
			previousModuleWord = "modules"
		}

		r.ReportIssue(fmt.Sprintf("The module count changed from the previous scan. It went from %d %s to %d.", previousModuleCount, previousModuleWord, currentModuleCount), report.IssueSeverityMedium)
		r.MakeRecommendation("Changes in the number of uploaded modules may be part of natural application lifecycle, but it is worth verifying.")
	}
}

func compareModuleSelection(r *report.Report, pr *report.Report) {
	currentSelectedModules := r.GetSelectedModules()
	previousSelectedModules := pr.GetSelectedModules()

	currentModuleNameCountMap := generateNameMappedArray(currentSelectedModules)
	previousModuleNameCountMap := generateNameMappedArray(previousSelectedModules)

	UNUSED(currentModuleNameCountMap, previousModuleNameCountMap)
}

func normalizeFilename(filename string) string {

	// Define a regular expression to match version numbers and "SNAPSHOT"
	// It ensures that the version number is at the end of the filename, preceded by a hyphen,
	// and followed by a file extension.
	re, _ := regexp.Compile(`(-\d+(\.\d+)*(-SNAPSHOT)?)((\.[a-zA-Z0-9]+)+)$`)

	// Replace version numbers and "SNAPSHOT" with an empty string
	normalized := re.ReplaceAllString(filename, "$4")

	return normalized
}

type ModuleNameCount struct {
	Counter int
	Strings []string
}

// generateNameMappedArray takes a slice of modules, normalizes the name
// and returns a map of module names to the number of times they appear in the slice
// and a slice of the original module names.
func generateNameMappedArray(modules []report.Module) map[string]ModuleNameCount {
	moduleNameCountMap := make(map[string]ModuleNameCount)

	for _, module := range modules {
		normalizedName := normalizeFilename(module.Name)

		if _, exists := moduleNameCountMap[normalizedName]; exists {
			moduleNameCount := moduleNameCountMap[normalizedName]
			moduleNameCount.Counter++
			moduleNameCount.Strings = append(moduleNameCount.Strings, normalizedName)
			moduleNameCountMap[normalizedName] = moduleNameCount
		} else {
			// If it doesn't exist, create a new entry
			moduleNameCount := ModuleNameCount{
				Counter: 1, // Start with 1 because it's the first occurrence
				Strings: []string{module.Name},
			}
			moduleNameCountMap[normalizedName] = moduleNameCount
		}
	}

	return moduleNameCountMap
}

func previousScan(r *report.Report, pr *report.Report) {

	if pr.Scan.BuildId == 0 {
		return
	}

	// SubScan types
	// Modules:
	// Compare module number - different number of modules > DONE.
	// Compare modules themselves - how do we distinguish versions?
	//   - if we have test-1.1.jar and test-2.0.jar in one scan, does that count as 2 instances of test.jar?
	// Compare module selection - different modules selected (needs to account for version numbers)
	// Compare scan size (if calculable) - warning not error if scan size has changed by, say, 20%
	// Compare technologies? Or is compare module selection enough for customers who change what they're scanning?
	// Files:
	// Compare number of files uploaded?

	compareModuleCount(r, pr)
	// compareModuleSelection(r, pr)

	return

	// TODO

	r.ReportIssue("The uploaded modules for this scan do not match the modules you uploaded for the previous scan. In this scan TODO modules were identified, and TODO were selected for scanning, whereas in the previous scan we observed TODO modules, TODO of which had been selected for scanning. Also noticeable was the total analysis size difference between the two scans.", report.IssueSeverityMedium)
	r.MakeRecommendation("The use of automation will lead to consistent scans.")
}
