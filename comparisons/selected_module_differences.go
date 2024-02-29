package comparisons

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"strconv"
	"strings"
)

func reportSelectedModuleDifferences(a, b *report.Report) {
	var r strings.Builder

	compareSelectedModuleDifferences(&r, "A", a, b)
	compareSelectedModuleDifferences(&r, "B", a, b)

	if r.Len() > 0 {
		utils.PrintTitle("Differences of Modules Selected As An Entry Point For Scanning")
		utils.ColorPrintf(r.String())
	}
}

func compareSelectedModuleDifferences(r *strings.Builder, side string, a, b *report.Report) {
	reportForThisSide := a
	reportForOtherSide := b

	if side == "B" {
		reportForThisSide = b
		reportForOtherSide = a
	}

	for _, moduleFoundInThisSide := range reportForThisSide.GetSelectedModules() {
		// We only care for module differences
		if isModuleInOtherReportSelectedModules(moduleFoundInThisSide, reportForOtherSide) {
			continue
		}

		var formattedSupportIssues = ""
		var formattedMissingSupportedFiles = ""
		var formattedSize = ""
		var formattedMd5 = ""
		var formattedPlatform = ""
		var isDependency = ""

		// Note We are only working with the first instance we find, this needs some thinking, should we or work across all + any differences for duplicates
		for _, instance := range moduleFoundInThisSide.Instances {
			if instance.Issues != nil {
				formattedSupportIssues = fmt.Sprintf(", %s", color.HiYellowString("Support issues = %d", len(instance.Issues)))
			}

			if len(formattedMissingSupportedFiles) < 1 && getMissingSupportedFileCount(instance) > 0 {
				formattedMissingSupportedFiles = fmt.Sprintf(", %s", color.HiYellowString("Missing Supporting Files = %d", getMissingSupportedFileCount(instance)))
			}

			if len(formattedSize) < 1 && len(instance.Size) > 0 {
				formattedSize = fmt.Sprintf(", Size = %s", instance.Size)
			}

			if len(formattedMd5) < 1 && len(instance.MD5) > 0 {
				formattedMd5 = fmt.Sprintf(", MD5 = %s", instance.MD5)
			}

			if len(formattedPlatform) < 1 && len(instance.Architecture) > 0 {
				formattedPlatform = fmt.Sprintf(", Platform = %s/%s/%s", instance.Architecture, instance.OperatingSystem, instance.Compiler)
			}

			if len(isDependency) < 1 && instance.IsDependency {
				isDependency = "Module is Dependency"
			}
		}

		r.WriteString(fmt.Sprintf("%s: \"%s\"%s%s%s%s%s%s\n",
			utils.GetFormattedOnlyInSideString(side),
			moduleFoundInThisSide.Name,
			formattedSize,
			formattedSupportIssues,
			formattedMissingSupportedFiles,
			isDependency,
			formattedMd5,
			formattedPlatform))
	}
}

func isModuleInOtherReportSelectedModules(module report.Module, other *report.Report) bool {
	for _, otherSelectedModule := range other.GetSelectedModules() {
		if module.Name == otherSelectedModule.Name {
			return true
		}
	}

	return false
}

func getMissingSupportedFileCount(instance report.ModuleInstance) int {
	for _, issue := range instance.Issues {
		if strings.HasPrefix(issue, "Missing Supporting Files") {
			trimmedPrefix := strings.Replace(issue, "Missing Supporting Files - ", "", 1)
			count, err := strconv.Atoi(strings.Split(trimmedPrefix, " ")[0])

			if err != nil {
				return 0
			}

			return count
		}
	}

	return 0
}
