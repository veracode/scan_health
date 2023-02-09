package main

import (
	"sync"
)

type Data struct {
	DetailedReport    DetailedReport
	PrescanFileList   PrescanFileList
	PrescanModuleList PrescanModuleList
}

func (api API) getData(appId, buildId int) Data {
	var data = Data{}

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

	data.DetailedReport.SubmittedDate = parseVeracodeDate(data.DetailedReport.StaticAnalysis.SubmittedDate).Local()
	data.DetailedReport.PublishedDate = parseVeracodeDate(data.DetailedReport.StaticAnalysis.PublishedDate).Local()
	data.DetailedReport.Duration = data.DetailedReport.PublishedDate.Sub(data.DetailedReport.SubmittedDate)
	data.DetailedReport.Duration = data.DetailedReport.PublishedDate.Sub(data.DetailedReport.SubmittedDate)

	return data
}
