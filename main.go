package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func main() {
	fmt.Printf("Scan Health v%s\nCopyright © Veracode, Inc. 2023. All Rights Reserved.\nThis is an unofficial Veracode product. It does not come with any support or warrenty.\n\n", AppVersion)
	vid := flag.String("vid", "", "Veracode API ID - See https://docs.veracode.com/r/t_create_api_creds")
	vkey := flag.String("vkey", "", "Veracode API key - See https://docs.veracode.com/r/t_create_api_creds")
	profile := flag.String("profile", "default", "Veracode credential profile - See https://docs.veracode.com/r/c_httpie_tool")
	region := flag.String("region", "", "Veracode Region [commercial, us, european]")
	scan := flag.String("sast", "", "Veracode Platform URL or build ID for a SAST application health review")

	flag.Parse()

	if !(*region == "" || *region == "commercial" || *region == "us" || *region == "european") {
		color.HiRed("Error: Invalid region. Must be either \"glocommercialbal\", \"us\" or \"european\"")
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
	data.reportDuplicateFiles()
	data.analyzeModules()
	data.analyzeModuleFatalErrors()
	data.analyzeModuleWarnings()
	data.outputRecommendations(api.region)
}

func (data Data) assertPrescanModulesPresent() {
	if len(data.PrescanModuleList.Modules) == 0 {
		color.HiRed("Error: Could not retrieve pre-scan modules")
		os.Exit(1)
	}
}
