package main

import (
	"fmt"
	"github.com/veracode/scan_health/v2/data"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"sync"
)

func performSASTCompare(scanA, scanB *string, api data.API, regionToUse string, outputFormat, jsonFilePath *string) {
	var wg sync.WaitGroup
	var scanABuildId int
	var scanBBuildId int

	wg.Add(2)

	go func() {
		defer wg.Done()
		buildId, err := utils.ParseBuildIdFromScanInformation(*scanA)
		if err != nil {
			utils.ErrorAndExit(fmt.Sprintf("Could not parse build ID from %s", *scanA), err)
		}
		scanABuildId = buildId
	}()

	go func() {
		defer wg.Done()
		buildId, err := utils.ParseBuildIdFromScanInformation(*scanB)
		if err != nil {
			utils.ErrorAndExit(fmt.Sprintf("Could not parse build ID from %s", *scanB), err)
		}
		scanBBuildId = buildId
	}()

	// Wait for the build IDs to load
	wg.Wait()

	if scanABuildId == scanBBuildId {
		utils.ErrorAndExit("These are both the same scan", nil)
	}

	utils.ColorPrintf(fmt.Sprintf("Comparing SAST scan build ID %d against %d in the %s region\n",
		scanABuildId,
		scanBBuildId,
		regionToUse))

	wg.Add(2)

	scanAReport := report.NewReport(scanABuildId, regionToUse, AppVersion, false)
	scanBReport := report.NewReport(scanBBuildId, regionToUse, AppVersion, false)

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

	print("TODO")

}
