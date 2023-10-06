package data

import (
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
	"sync"
)

func (api API) PopulateReportWithDataFromAPI(report *report.Report) {
	var wg sync.WaitGroup

	executeApiCall(report, &wg, api.populateDetailedReport)

	// Wait for the detailed report
	wg.Wait()

	executeApiCall(report, &wg, api.populateAppInfo)

	// Wait for the app info
	wg.Wait()

	if report.Scan.ApplicationId == 0 {
		utils.ErrorAndExit("Could not resolve the application ID because only a build ID was supplied, and the scan has not finished. Please try again using the URL instead of the build ID.", nil)
	}

	if report.Scan.ScanName == "" {
		executeApiCall(report, &wg, api.populateBuildInfo)
	}

	executeApiCall(report, &wg, api.populateSandboxInfo)
	executeApiCall(report, &wg, api.getPrescanFileList)
	executeApiCall(report, &wg, api.getPrescanModuleList)
	executeApiCall(report, &wg, api.populatePreviousBuildInfo)

	// Wait for all the data to load
	wg.Wait()

	postProcessData(report)
}

func executeApiCall(report *report.Report, wg *sync.WaitGroup, fn func(report *report.Report)) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		fn(report)
	}()
}
