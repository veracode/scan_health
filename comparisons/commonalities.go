package comparisons

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"strings"
)

func reportCommonalities(a, b *report.Report) {
	var report strings.Builder

	if a.Scan.ApplicationName == b.Scan.ApplicationName {
		report.WriteString(fmt.Sprintf("Application:        \"%s\"\n", a.Scan.ApplicationName))
	}

	if a.Scan.SandboxId == b.Scan.SandboxId && len(a.Scan.SandboxName) > 0 {
		report.WriteString(fmt.Sprintf("Sandbox:            \"%s\"\n", a.Scan.SandboxName))
	}

	if a.Scan.ScanName == b.Scan.ScanName {
		report.WriteString(fmt.Sprintf("Scan name:          \"%s\"\n", a.Scan.ScanName))
	}

	if len(a.UploadedFiles) == len(b.UploadedFiles) {
		report.WriteString(fmt.Sprintf("Files uploaded:     %d\n", len(a.UploadedFiles)))
	}

	// TODO: check this is the correct count - we infer moduels from flaws.
	if len(a.Modules) == len(b.Modules) {
		report.WriteString(fmt.Sprintf("Total modules:      %d\n", len(a.Modules)))
	}

	if len(a.GetSelectedModules()) == len(b.GetSelectedModules()) {
		report.WriteString(fmt.Sprintf("Modules selected:   %d\n", len(a.GetSelectedModules())))
	}

	if a.Scan.EngineVersion == b.Scan.EngineVersion {
		report.WriteString(fmt.Sprintf("Engine version:     %s\n", a.Scan.EngineVersion))
	}

	if a.Flaws.Total == b.Flaws.Total && a.Flaws.Mitigated == b.Flaws.Mitigated && a.Flaws.TotalAffectingPolicy == b.Flaws.TotalAffectingPolicy && a.Flaws.OpenAffectingPolicy == b.Flaws.OpenAffectingPolicy && a.Flaws.OpenButNotAffectingPolicy == b.Flaws.OpenButNotAffectingPolicy {
		flawsFormatted := fmt.Sprintf("Flaws:              %d total, %d mitigated, %d policy affecting, %d open affecting policy, %d open not affecting policy\n", a.Flaws.Total, a.Flaws.Mitigated, a.Flaws.TotalAffectingPolicy, a.Flaws.OpenAffectingPolicy, a.Flaws.OpenButNotAffectingPolicy)

		if a.Flaws.Total == 0 {
			report.WriteString(color.HiYellowString(flawsFormatted))
		} else {
			report.WriteString(flawsFormatted)
		}
	}

	if report.Len() > 0 {
		utils.PrintTitle("In common with both scans")
		utils.ColorPrintf(report.String())
	}
}
