package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/veracode/scan_health/v2/checks"
	"github.com/veracode/scan_health/v2/data"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"os"
)

func performSASTHealthCheck(scan *string, api data.API, regionToUse string, includePreviousScan *bool, outputFormat, jsonFilePath *string, errorOnHighSeverity *bool) {
	buildId, err := utils.ParseBuildIdFromScanInformation(*scan)
	if err != nil {
		utils.ErrorAndExit("Could not parse build ID", err)
	}

	utils.ColorPrintf(fmt.Sprintf("Inspecting SAST build id = %d in the %s region\n",
		buildId,
		regionToUse))

	healthReport := report.NewReport(buildId, regionToUse, AppVersion, false)

	// Try to set the application ID however we could be working from a build ID so this may not be available
	applicationId, err := utils.ParseApplicationIdFromPlatformUrl(*scan)

	if err == nil {
		healthReport.Scan.ApplicationId = applicationId
	}

	api.PopulateReportWithDataFromAPI(healthReport, *includePreviousScan)

	if !healthReport.Scan.IsLatestScan {
		if len(healthReport.Scan.SandboxName) > 0 {
			color.HiYellow("Warning: This is not the latest SAST scan in this sandbox")
		} else {
			color.HiYellow("Warning: This is not the latest SAST policy scan")
		}
	}

	var previousHealthReport = &report.Report{}

	if *includePreviousScan == true {
		previousBuildId := api.GetPreviousBuildId(healthReport)

		if previousBuildId > 0 {
			previousHealthReport = report.NewReport(previousBuildId, regionToUse, AppVersion, true)
			previousHealthReport.Scan.ApplicationId = healthReport.Scan.ApplicationId
			api.PopulateReportWithDataFromAPI(previousHealthReport, false)
		}
	}

	checks.PerformChecks(healthReport, previousHealthReport)

	healthReport.PrioritizeIssues()

	healthReport.Render(*outputFormat, *jsonFilePath)

	if *errorOnHighSeverity {
		// Return exit code of 1 if any high severity issues found
		for _, issue := range healthReport.Issues {
			if issue.Severity == report.IssueSeverityHigh {
				os.Exit(1)
			}
		}
	}
}
