package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/veracode/scan_health/v2/checks"
	"github.com/veracode/scan_health/v2/data"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
)

func main() {
	fmt.Printf("Scan Health v%s\nCopyright © Veracode, Inc. 2023. All Rights Reserved.\nThis is an unofficial Veracode product. It does not come with any support or warranty.\n\n", AppVersion)
	vid := flag.String("vid", "", "Veracode API ID - See https://docs.veracode.com/r/t_create_api_creds")
	vkey := flag.String("vkey", "", "Veracode API key - See https://docs.veracode.com/r/t_create_api_creds")
	profile := flag.String("profile", "default", "Veracode credential profile - See https://docs.veracode.com/r/c_httpie_tool#using-multiple-profiles")
	region := flag.String("region", "", "Veracode Region [commercial, us, european]. Required if a Build ID is specified.")
	scan := flag.String("sast", "", "Veracode Platform URL or build ID for a SAST application health review")
	outputFormat := flag.String("format", "console", "Output format [console, json]")
	jsonFilePath := flag.String("json-file", "", "Optional file for writing JSON output to")
	includePreviousScan := flag.Bool("previous-scan", false, "Enable comparison with the previous scan (this will result in many requests being made)")
	enableCaching := flag.Bool("cache", false, "Enable caching of API responses (useful for development)")
	errorOnHighSeverity := flag.Bool("error-on-high-severity", false, "Return a non-zero exit code if any high severity issues are found")

	flag.Parse()

	if *region != "" && utils.IsValidRegion(*region) == false {
		utils.ErrorAndExitWithUsage(fmt.Sprintf("Invalid region \"%s\". Must be either \"commercial\", \"us\" or \"european\"", *region))
	}

	if *region != "" &&
		(strings.HasPrefix(*scan, "https://") && utils.ParseRegionFromUrl(*scan) != *region) {
		utils.ErrorAndExit(fmt.Sprintf("The region from the URL (%s) does not match that specified by the command line (%s)", utils.ParseRegionFromUrl(*scan), *region), nil)
	}

	if len(*scan) < 1 {
		utils.ErrorAndExitWithUsage("No Veracode Platform URL or build ID specified for the health review. Expected: \"scan_health -sast https://analysiscenter.veracode.com/auth/index.jsp...\"")
	}

	if !(*outputFormat == "console" || *outputFormat == "json") {
		utils.ErrorAndExitWithUsage("Invalid output format. Must be either \"console\"or \"json\"")
	}

	var regionToUse string

	// Command line region (if specified) takes precedence
	if *region != "" {
		regionToUse = *region
	} else {
		regionToUse = utils.ParseRegionFromUrl(*scan)
	}

	notifyOfUpdates()

	apiId, apiKey := getCredentials(*vid, *vkey, *profile)
	api := data.API{Id: apiId, Key: apiKey, Region: regionToUse, AppVersion: AppVersion, EnableCaching: *enableCaching}

	buildId, err := utils.ParseBuildIdFromScanInformation(*scan)
	if err != nil {
		utils.ErrorAndExit("", err)
	}

	api.AssertCredentialsWork()

	utils.ColorPrintf(fmt.Sprintf("Inspecting SAST build id = %d in the %s region\n",
		buildId,
		regionToUse))

	healthReport := report.NewReport(buildId, regionToUse, AppVersion, false)
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
