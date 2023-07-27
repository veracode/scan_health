package main

import (
	"flag"
	"fmt"
	"github.com/antfie/scan_health/v2/checks"
	"github.com/antfie/scan_health/v2/data"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
	"strings"
)

func main() {
	fmt.Printf("Scan Health v%s\nCopyright © Veracode, Inc. 2023. All Rights Reserved.\nThis is an unofficial Veracode product. It does not come with any support or warranty.\n\n", AppVersion)
	vid := flag.String("vid", "", "Veracode API ID - See https://docs.veracode.com/r/t_create_api_creds")
	vkey := flag.String("vkey", "", "Veracode API key - See https://docs.veracode.com/r/t_create_api_creds")
	profile := flag.String("profile", "default", "Veracode credential profile - See https://docs.veracode.com/r/c_httpie_tool#using-multiple-profiles")
	region := flag.String("region", "", "Veracode Region [commercial, us, european]")
	scan := flag.String("sast", "", "Veracode Platform URL or build ID for a SAST application health review")
	outputFormat := flag.String("format", "console", "Output format [console, json]")
	jsonFilePath := flag.String("json-file", "", "Optional file for writing JSON output to")
	enableCaching := flag.Bool("cache", false, "Enable caching")

	flag.Parse()

	if !(*region == "" || *region == "commercial" || *region == "us" || *region == "european") {
		utils.ErrorAndExitWithUsage("Invalid region. Must be either \"commercial\", \"us\" or \"european\"")
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
	buildId := utils.ParseBuildIdFromPlatformUrl(*scan)

	api.AssertCredentialsWork()

	utils.ColorPrintf(fmt.Sprintf("Inspecting SAST build id = %d in the %s region\n",
		buildId,
		regionToUse))

	healthReport := report.NewReport(buildId, regionToUse, AppVersion)

	api.PopulateReportWithDataFromAPI(healthReport)

	checks.PerformChecks(healthReport)

	healthReport.PrioritizeIssues()

	healthReport.Render(*outputFormat, *jsonFilePath)
}
