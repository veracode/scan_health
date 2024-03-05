package comparisons

import (
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

		reportOnModules(r, side, moduleFoundInThisSide)
	}
}

func isModuleInOtherReportSelectedModules(module report.Module, other *report.Report) bool {
	return module.IsInListByName(other.GetSelectedModules())
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
