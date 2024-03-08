package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/veracode/scan_health/v2/data"
	"github.com/veracode/scan_health/v2/utils"
)

func main() {
	fmt.Printf("Scan Health v%s\nCopyright © Veracode, Inc. 2024. All Rights Reserved.\nThis is an unofficial Veracode product. It does not come with any support or warranty.\n\n", AppVersion)
	vid := flag.String("vid", "", "Veracode API ID - See https://docs.veracode.com/r/c_api_credentials3")
	vkey := flag.String("vkey", "", "Veracode API key - See https://docs.veracode.com/r/c_api_credentials3")
	profile := flag.String("profile", "default", "Veracode credential profile - See https://docs.veracode.com/r/c_httpie_tool#using-multiple-profiles")
	region := flag.String("region", "", "Veracode Region [commercial, us, european], required if a Build ID is specified")
	action := flag.String("action", "health", "An action to perform [health, compare]")
	scan := flag.String("sast", "", "Veracode Platform URL or build ID for a SAST scan")
	scanA := flag.String("a", "", "Veracode Platform URL or build ID for SAST scan \"A\" - for scan comparison")
	scanB := flag.String("b", "", "Veracode Platform URL or build ID for SAST scan \"B\" - for scan comparison")
	outputFormat := flag.String("format", "console", "Output format [console, json]")
	jsonFilePath := flag.String("json-file", "", "Optional file for writing JSON output to")
	includePreviousScan := flag.Bool("previous-scan", false, "Enable comparison with the previous scan (this will result in many requests being made)")
	enableCaching := flag.Bool("cache", false, "Enable caching of API responses (useful for development). It is not recommended to use caching as sensitive data will be stored on disk.\n")
	errorOnHighSeverity := flag.Bool("error-on-high-severity", false, "Return a non-zero exit code if any high severity issues are found")

	flag.Parse()

	if *action != "health" && *action != "compare" {
		utils.ErrorAndExitWithUsage(fmt.Sprintf("Invalid action \"%s\". Must be either \"health\" or \"compare\"", *action))
	}

	if *region != "" && utils.IsValidRegion(*region) == false {
		utils.ErrorAndExitWithUsage(fmt.Sprintf("Invalid region \"%s\". Must be either \"commercial\", \"us\" or \"european\"", *region))
	}

	if *region != "" &&
		(strings.HasPrefix(*scan, "https://") && utils.ParseRegionFromUrl(*scan) != *region) {
		utils.ErrorAndExit(fmt.Sprintf("The region from the URL (%s) does not match that specified by the command line (%s)", utils.ParseRegionFromUrl(*scan), *region), nil)
	}

	var regionToUse string

	// Region from the command line takes precedence
	if len(*region) > 0 {
		regionToUse = *region
	} else {
		regionToUse = utils.ParseRegionFromUrl(*scan)
	}

	notifyOfUpdates()

	apiId, apiKey := getCredentials(*vid, *vkey, *profile)
	api := data.API{Id: apiId, Key: apiKey, Region: regionToUse, AppVersion: AppVersion, EnableCaching: *enableCaching}

	api.AssertCredentialsWork()

	switch *action {
	case "health":
		if len(*scan) < 1 {
			utils.ErrorAndExitWithUsage("No Veracode Platform URL or build ID specified for the health review. Expected: \"scan_health -sast https://analysiscenter.veracode.com/auth/index.jsp...\"")
		}

		if !(*outputFormat == "console" || *outputFormat == "json") {
			utils.ErrorAndExitWithUsage("Invalid output format. Must be either \"console\"or \"json\"")
		}

		performSASTHealthCheck(scan, api, regionToUse, includePreviousScan, outputFormat, jsonFilePath, errorOnHighSeverity)
	case "compare":
		if len(*scanA) < 1 {
			utils.ErrorAndExitWithUsage("No Veracode Platform URL or build ID specified for SAST scan A. Expected: \"scan_health -action compare -a https://analysiscenter.veracode.com/auth/index.jsp... -b https://analysiscenter.veracode.com/auth/index.jsp...\"")
		}

		if len(*scanB) < 1 {
			utils.ErrorAndExitWithUsage("No Veracode Platform URL or build ID specified for SAST scan B. Expected: \"scan_health -action compare -a https://analysiscenter.veracode.com/auth/index.jsp... -b https://analysiscenter.veracode.com/auth/index.jsp...\"")
		}

		if utils.ParseRegionFromUrl(*scanA) != utils.ParseRegionFromUrl(*scanB) {
			utils.ErrorAndExit("The SAST scans are for different regions. We can only compare SAST scans within the same region.", nil)
		}

		performSASTCompare(scanA, scanB, api, regionToUse)
	}
}
