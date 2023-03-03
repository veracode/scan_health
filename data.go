package main

import (
	"os"
	"sync"

	"github.com/fatih/color"
)

type Data struct {
	DetailedReport    DetailedReport
	PrescanFileList   PrescanFileList
	PrescanModuleList PrescanModuleList
	Recommendations   *[]string
}

func (data Data) makeRecommendation(recommendation string) {
	if !isStringInStringArray(recommendation, *data.Recommendations) {
		*data.Recommendations = append(*data.Recommendations, recommendation)
	}
}

func (api API) getData(appId, buildId int) Data {
	var data = Data{Recommendations: &[]string{}}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		data.DetailedReport = api.getDetailedReport(buildId)
	}()

	wg.Wait()

	// We can't rely on the passed-in app IDs as they may not be present (if not using a URL).
	// Try to get the app ID from the detailed report if possible
	if data.DetailedReport.DataAvailable {
		appId = data.DetailedReport.AppId
	}

	if appId == -1 {
		color.HiRed("Error: Could not resolve the application ID because only a build ID was supplied, and the scan has not finished. Please try again using the URL instead of the build ID.")
		os.Exit(1)
	}

	// Resolve other information if we cannot get the detailed report
	if !data.DetailedReport.DataAvailable {
		wg.Add(2)

		go func() {
			defer wg.Done()
			api.populateReportDetailsFromAppInfo(appId, &data.DetailedReport)
		}()

		go func() {
			defer wg.Done()
			api.populateReportDetailsFromBuildInfo(appId, buildId, &data.DetailedReport)
		}()
	}

	wg.Add(2)

	go func() {
		defer wg.Done()
		data.PrescanFileList = api.getPrescanFileList(appId, buildId)
	}()

	go func() {
		defer wg.Done()
		data.PrescanModuleList = api.getPrescanModuleList(appId, buildId)
	}()

	wg.Wait()

	return data
}
