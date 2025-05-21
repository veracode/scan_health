package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/veracode/scan_health/v2/comparisons"
	"github.com/veracode/scan_health/v2/data"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"sync"
)

func performSASTCompare(scanA, scanB *string, api data.API, regionToUse string) {
	var wg sync.WaitGroup
	var scanABuildId int
	var scanBBuildId int

	wg.Add(2)

	go func() {
		defer wg.Done()

		buildId, err := utils.ParseBuildIdFromScanInformation(*scanA)
		if err != nil {
			utils.ErrorAndExit(fmt.Sprintf("Could not parse build ID from: \"%s\"", *scanA), err)
		}

		scanABuildId = buildId
	}()

	go func() {
		defer wg.Done()

		buildId, err := utils.ParseBuildIdFromScanInformation(*scanB)
		if err != nil {
			utils.ErrorAndExit(fmt.Sprintf("Could not parse build ID from: \"%s\"", *scanB), err)
		}

		scanBBuildId = buildId
	}()

	// Wait for the build IDs to load
	wg.Wait()

	if scanABuildId == scanBBuildId {
		utils.ErrorAndExit("These are both the same scan", nil)
	}

	utils.ColorPrintf(fmt.Sprintf("Comparing SAST scan %s against scan %s in the %s region\n",
		color.HiGreenString("\"A\" (Build id = %d)", scanABuildId),
		color.HiMagentaString("\"B\" (Build id = %d)", scanBBuildId),
		regionToUse))

	wg.Add(2)

	scanAReport := report.NewReport(scanABuildId, regionToUse, AppVersion, false)
	scanBReport := report.NewReport(scanBBuildId, regionToUse, AppVersion, false)

	// Try to set the application ID however we could be working from a build ID so this may not be available
	scanAApplicationId, err := utils.ParseApplicationIdFromPlatformUrl(*scanA)

	if err == nil {
		scanAReport.Scan.ApplicationId = scanAApplicationId
	}

	// Try to set the application ID however we could be working from a build ID so this may not be available
	scanBApplicationId, err := utils.ParseApplicationIdFromPlatformUrl(*scanB)

	if err == nil {
		scanBReport.Scan.ApplicationId = scanBApplicationId
	}

	go func() {
		defer wg.Done()
		api.PopulateReportWithDataFromAPI(scanAReport, false)
	}()

	go func() {
		defer wg.Done()
		api.PopulateReportWithDataFromAPI(scanBReport, false)
	}()

	// Wait for all the data to load
	wg.Wait()

	comparisons.PerformComparisons(scanAReport, scanBReport)
}
