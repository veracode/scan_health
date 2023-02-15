package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func (data Data) reportScanDetails(region string) {
	println()
	printTitle("Scan Summary")

	fmt.Printf("Account ID:         %d\n", data.DetailedReport.AccountId)
	fmt.Printf("Application:        \"%s\"\n", data.DetailedReport.AppName)

	if len(data.DetailedReport.SandboxName) > 0 {
		fmt.Printf("Sandbox:            \"%s\"\n", data.DetailedReport.SandboxName)
	}

	fmt.Printf("Scan name:          \"%s\"\n", data.DetailedReport.StaticAnalysis.ScanName)
	fmt.Printf("Review Modules URL: %s\n", data.DetailedReport.getReviewModulesUrl(region))
	fmt.Printf("Triage Flaws URL:   %s\n", data.DetailedReport.getTriageFlawsUrl(region))
	fmt.Printf("Files uploaded:     %d\n", len(data.PrescanFileList.Files))
	fmt.Printf("Total modules:      %d\n", len(data.PrescanModuleList.Modules))
	fmt.Printf("Modules selected:   %d\n", len(data.DetailedReport.StaticAnalysis.Modules))
	fmt.Printf("Engine version:     %s\n", data.DetailedReport.StaticAnalysis.EngineVersion)
	fmt.Printf("Submitted:          %s (%s ago)\n", data.DetailedReport.SubmittedDate, formatDuration(time.Since(data.DetailedReport.SubmittedDate)))
	fmt.Printf("Published:          %s (%s ago)\n", data.DetailedReport.PublishedDate, formatDuration(time.Since(data.DetailedReport.PublishedDate)))
	fmt.Printf("Duration:           %s\n", data.DetailedReport.Duration)

	flawsFormatted := fmt.Sprintf("Flaws:              %d total, %d mitigated, %d policy affecting, %d open affecting policy, %d open not affecting policy\n", data.DetailedReport.TotalFlaws, data.DetailedReport.TotalFlaws-data.DetailedReport.UnmitigatedFlaws, data.DetailedReport.getPolicyAffectingFlawCount(), data.DetailedReport.getOpenPolicyAffectingFlawCount(), data.DetailedReport.getOpenNonPolicyAffectingFlawCount())

	if data.DetailedReport.TotalFlaws == 0 {
		color.HiYellow(flawsFormatted)
		data.makeRecommendation("When no flaws have been found this can be an indication that the scan is misconfigured")
	} else {
		fmt.Print(flawsFormatted)
	}

	println()
}
