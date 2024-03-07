package comparisons

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"strings"
)

func reportNotSelectedModuleDifferences(a, b *report.Report) {
	var r strings.Builder

	compareTopLevelNotSelectedModules(&r, "A", a, b, false)
	compareTopLevelNotSelectedModules(&r, "B", a, b, false)

	if r.Len() > 0 {
		if strings.Contains(r.String(), "files extracted from") {
			utils.PrintTitle("Differences of Top-Level Modules Which May or May Not Have Been Selected")
		} else {
			utils.PrintTitle("Differences of Top-Level Modules Not Selected As An Entry Point (And Not Scanned) - Unselected Potential First Party Components")
		}

		utils.ColorPrintf(r.String())
	}
}

func reportDependencyModuleDifferences(a, b *report.Report) {
	var r strings.Builder

	compareTopLevelNotSelectedModules(&r, "A", a, b, true)
	compareTopLevelNotSelectedModules(&r, "B", a, b, true)

	if r.Len() > 0 {
		utils.PrintTitle("Differences of Dependency Modules Not Selected As An Entry Point")
		utils.ColorPrintf(r.String())
	}
}

func compareTopLevelNotSelectedModules(r *strings.Builder, side string, a, b *report.Report, onlyDependencies bool) {
	reportForThisSide := a
	reportForOtherSide := b

	if side == "B" {
		reportForThisSide = b
		reportForOtherSide = a
	}

	for _, moduleFoundInThisSide := range reportForThisSide.GetPrescanModules() {
		if moduleFoundInThisSide.IsDependency() != onlyDependencies {
			continue
		}

		// Ignore selected modules
		if moduleFoundInThisSide.IsSelected() {
			continue
		}

		// We only care for module differences where the module is unselected and present in both sides
		if isModuleInOtherReportUnselectedModules(moduleFoundInThisSide, reportForOtherSide) {
			continue
		}

		reportOnModules(r, side, moduleFoundInThisSide)
	}
}

func reportOnModules(r *strings.Builder, side string, moduleFoundInThisSide report.Module) {
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

		if len(formattedPlatform) < 1 {
			if len(instance.Architecture) > 0 {
				formattedPlatform = fmt.Sprintf(", Platform = %s/%s/%s", instance.Architecture, instance.OperatingSystem, instance.Compiler)
			} else if len(instance.Platform) > 0 {
				formattedPlatform = fmt.Sprintf(", Platform = %s", instance.Platform)
			}
		}

		if len(isDependency) < 1 && instance.IsDependency {
			isDependency = " Module is Dependency"
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

func isModuleInOtherReportUnselectedModules(module report.Module, other *report.Report) bool {
	for _, otherSelectedModule := range other.Modules {
		// Ignore cases when both are selected
		if module.IsSelected() && otherSelectedModule.IsSelected() {
			return false
		}

		if module.Name == otherSelectedModule.Name {
			return true
		}
	}

	return false
}
