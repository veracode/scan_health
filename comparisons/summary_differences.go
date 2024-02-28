package comparisons

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"time"
)

func reportSummaryDifferences(side string, a, b *report.Report) {
	utils.ColorPrintf(getFormattedSideStringWithMessage(side, fmt.Sprintf("\nScan %s", side)))
	fmt.Println("\n======")

	reportForThisSide := a

	if side == "B" {
		reportForThisSide = b
	}

	if a.Scan.AccountId != b.Scan.AccountId {
		fmt.Printf("Account ID:         %d\n", reportForThisSide.Scan.AccountId)
	}

	if a.Scan.ApplicationName != b.Scan.ApplicationName {
		fmt.Printf("Application:        \"%s\"\n", reportForThisSide.Scan.ApplicationName)
	}

	if a.Scan.SandboxId != b.Scan.SandboxId && len(reportForThisSide.Scan.SandboxName) > 0 {
		fmt.Printf("Sandbox:            \"%s\"\n", reportForThisSide.Scan.SandboxName)
	}

	if a.Scan.ScanName != b.Scan.ScanName {
		fmt.Printf("Scan name:          \"%s\"\n", reportForThisSide.Scan.ScanName)
	}

	// TODO: fix up these URLs
	//fmt.Printf("Review Modules URL: %s\n", thisDetailedReport.getReviewModulesUrl(region))
	//fmt.Printf("Triage Flaws URL:   %s\n", thisDetailedReport.getTriageFlawsUrl(region))

	if len(a.UploadedFiles) != len(b.UploadedFiles) {
		fmt.Printf("Files uploaded:     %d\n", len(reportForThisSide.UploadedFiles))
	}

	if len(a.Modules) != len(b.Modules) {
		fmt.Printf("Total modules:      %d\n", len(reportForThisSide.Modules))
	}

	if len(a.GetSelectedModules()) != len(b.GetSelectedModules()) {
		fmt.Printf("Modules selected:   %d\n", len(reportForThisSide.GetSelectedModules()))
	}

	if a.Scan.EngineVersion != b.Scan.EngineVersion {
		fmt.Printf("Engine version:     %s\n", reportForThisSide.Scan.EngineVersion)
	}

	fmt.Printf("Submitted:          %s (%s ago)\n", reportForThisSide.Scan.SubmittedDate, utils.FormatDuration(time.Since(reportForThisSide.Scan.SubmittedDate)))
	fmt.Printf("Published:          %s (%s ago)\n", reportForThisSide.Scan.PublishedDate, utils.FormatDuration(time.Since(reportForThisSide.Scan.PublishedDate)))
	fmt.Printf("Duration:           %s\n", utils.FormatDuration(reportForThisSide.Scan.ScanDuration))

	if !(a.Flaws.Total == b.Flaws.Total && a.Flaws.Mitigated == b.Flaws.Mitigated && a.Flaws.TotalAffectingPolicy == b.Flaws.TotalAffectingPolicy && a.Flaws.OpenAffectingPolicy == b.Flaws.OpenAffectingPolicy && a.Flaws.OpenButNotAffectingPolicy == b.Flaws.OpenButNotAffectingPolicy) {
		flawsFormatted := fmt.Sprintf("Flaws:              %d total, %d mitigated, %d policy affecting, %d open affecting policy, %d open not affecting policy\n", reportForThisSide.Flaws.Total, reportForThisSide.Flaws.Mitigated, reportForThisSide.Flaws.TotalAffectingPolicy, reportForThisSide.Flaws.OpenAffectingPolicy, reportForThisSide.Flaws.OpenButNotAffectingPolicy)

		if reportForThisSide.Flaws.Total == 0 {
			color.HiYellow(flawsFormatted)
		} else {
			fmt.Print(flawsFormatted)
		}
	}
}

func getFormattedSideStringWithMessage(side, message string) string {
	if side == "A" {
		return color.HiGreenString(message)
	}

	return color.HiMagentaString(message)
}
