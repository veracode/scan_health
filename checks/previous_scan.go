package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"regexp"
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

	currentModuleNameCountMap := generateMappedArray(currentSelectedModules)
	previousModuleNameCountMap := generateMappedArray(previousSelectedModules)

	// Now, how do we distinguish between selections?
	// What if one scan had 3 files called 'bob.jar' selected and another had 1?
	UNUSED(currentModuleNameCountMap, previousModuleNameCountMap)
}

func normalizeFilename(filename string) (string, error) {

	// Define a regular expression to match version numbers and "SNAPSHOT"
	// It ensures that the version number is at the end of the filename, preceded by a hyphen,
	// and followed by a file extension.
	re, err := regexp.Compile(`(-\d+(\.\d+)*(-SNAPSHOT)?)((\.[a-zA-Z0-9]+)+)$`)
	if err != nil {
		return "", err
	}

	// Replace version numbers and "SNAPSHOT" with an empty string
	normalized := re.ReplaceAllString(filename, "$4")

	return normalized, nil
}

func generateMappedArray(modules []report.Module) map[string]int {
	moduleNameCountMap := make(map[string]int)
	for _, module := range modules {
		normalizedName, err := normalizeFilename(module.Name)

		if err == nil {
			moduleNameCountMap[normalizedName]++
		} else {
			moduleNameCountMap[module.Name]++
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
	// Compare module selection - different modules selected (needs to account for version numbers)
	// Compare scan size (if calculable) - warning not error if scan size has changed by, say, 20%
	// Compare technologies? Or is compare module selection enough for customers who change what they're scanning?
	// Files:
	// Compare number of files uploaded?

	compareModuleCount(r, pr)
	compareModuleSelection(r, pr)

	return

	// TODO

	r.ReportIssue("The uploaded modules for this scan do not match the modules you uploaded for the previous scan. In this scan TODO modules were identified, and TODO were selected for scanning, whereas in the previous scan we observed TODO modules, TODO of which had been selected for scanning. Also noticeable was the total analysis size difference between the two scans.", report.IssueSeverityMedium)
	r.MakeRecommendation("The use of automation will lead to consistent scans.")
}
