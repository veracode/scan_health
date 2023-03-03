package main

import (
	"sync"
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
	var data = Data{}
	data.Recommendations = &[]string{}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		data.DetailedReport = api.getDetailedReport(buildId)
	}()

	wg.Wait()
	wg.Add(2)

	// We can't rely on the passed-in app IDs as they may not be present if not using a URL, so get the app ID from the detailed report

	go func() {
		defer wg.Done()
		data.PrescanFileList = api.getPrescanFileList(data.DetailedReport.AppId, buildId)
	}()

	go func() {
		defer wg.Done()
		data.PrescanModuleList = api.getPrescanModuleList(data.DetailedReport.AppId, buildId)
	}()

	wg.Wait()

	return data
}
