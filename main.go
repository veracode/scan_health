package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

func main() {
	fmt.Printf("Scan Health v%s\nCopyright © Veracode, Inc. 2023. All Rights Reserved.\nThis is an unofficial Veracode product. It does not come with any support or warrenty.\n\n", AppVersion)
	vid := flag.String("vid", "", "Veracode API ID - See https://docs.veracode.com/r/t_create_api_creds")
	vkey := flag.String("vkey", "", "Veracode API key - See https://docs.veracode.com/r/t_create_api_creds")
	profile := flag.String("profile", "default", "Veracode credential profile - See https://docs.veracode.com/r/c_httpie_tool")
	region := flag.String("region", "", "Veracode Region [global, us, eu]")
	scan := flag.String("sast", "", "Veracode Platform URL or build ID for a SAST application health review")

	flag.Parse()

	if !(*region == "" || *region == "global" || *region == "us" || *region == "eu") {
		color.HiRed("Error: Invalid region. Must be either \"global\", \"us\" or \"eu\"")
		print("\nUsage:\n")
		flag.PrintDefaults()
		return
	}

	if len(*scan) < 1 {
		color.HiRed("Error: No Veracode Platform URL or build ID specified for the health review. Expected: \"scan_health -sast https://analysiscenter.veracode.com/auth/index.jsp...\"")
		print("\nUsage:\n")
		flag.PrintDefaults()
		return
	}

	if *region != "" &&
		(strings.HasPrefix(*scan, "https://") && parseRegionFromUrl(*scan) != *region) {
		color.HiRed(fmt.Sprintf("Error: The region from the URL (%s) does not match that specified by the command line (%s)", parseRegionFromUrl(*scan), *region))
		os.Exit(1)
	}

	var regionToUse string

	// Command line region takes precidence
	if *region == "" {
		regionToUse = parseRegionFromUrl(*scan)
	} else {
		regionToUse = *region
	}

	notifyOfUpdates()

	var apiId, apiKey = getCredentials(*vid, *vkey, *profile)
	var api = API{apiId, apiKey, regionToUse}

	scanAppId := parseAppIdFromPlatformUrl(*scan)
	scanBuildId := parseBuildIdFromPlatformUrl(*scan)

	api.assertCredentialsWork()

	colorPrintf(fmt.Sprintf("Inspecting SAST build id = %d for health in the %s region\n",
		scanBuildId,
		api.region))

	data := api.getData(scanAppId, scanBuildId)

	data.reportScanDetails(api.region)
	data.assertPrescanModulesPresent()
	data.analyzeUploadedFiles()
}

func (data Data) reportScanDetails(region string) {
	fmt.Printf("Account ID:         %d\n", data.DetailedReport.AccountId)
	fmt.Printf("Application:        \"%s\"\n", data.DetailedReport.AppName)
	fmt.Printf("Sandbox:            \"%s\"\n", data.DetailedReport.SandboxName)
	fmt.Printf("Scan name:          \"%s\"\n", data.DetailedReport.StaticAnalysis.ScanName)
	fmt.Printf("Review Modules URL: %s\n", data.DetailedReport.getReviewModulesUrl(region))
	fmt.Printf("Traige Flaws URL:   %s\n", data.DetailedReport.getTriageFlawsUrl(region))
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
	} else {
		fmt.Print(flawsFormatted)
	}
}

func (data Data) assertPrescanModulesPresent() {
	if len(data.PrescanModuleList.Modules) == 0 {
		color.HiRed("Error: Could not retrieve pre-scan modules")
		os.Exit(1)
	}
}
