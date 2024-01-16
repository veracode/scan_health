package data

import (
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"sync"
)

func (api API) PopulateReportWithDataFromAPI(r *report.Report, includeDataFromPreviousScan bool) {
	var wg sync.WaitGroup

	executeApiCall(r, &wg, api.populateDetailedReport)

	// Wait for the detailed report
	wg.Wait()

	executeApiCall(r, &wg, api.populateAppInfo)

	// Unless we are looking at a previous scan..
	if !r.IsReportForOtherScan {
		if includeDataFromPreviousScan || len(r.Scan.ScanName) == 0 {
			// We need to make this call for the scan name if one was not retrieved from the Detailed Report
			// We need to call this if we are analysing previous scan
			executeApiCall(r, &wg, api.populateBuildList)
		}
	}

	// Wait for the app info and build list
	wg.Wait()

	if r.Scan.ApplicationId == 0 {
		utils.ErrorAndExit("Could not resolve the application ID because only a build ID was supplied, and the scan has not finished. Please try again using the URL instead of the build ID.", nil)
	}

	executeApiCall(r, &wg, api.populateSandboxInfo)
	executeApiCall(r, &wg, api.populatePrescanFileList)
	executeApiCall(r, &wg, api.populatePrescanModuleList)

	// Wait for all the data to load
	wg.Wait()

	postProcessData(r)
}

func executeApiCall(r *report.Report, wg *sync.WaitGroup, fn func(r *report.Report)) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		fn(r)
	}()
}
