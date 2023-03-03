package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type AppInfo struct {
	XMLName     xml.Name           `xml:"appinfo"`
	AccountId   int                `xml:"account_id,attr"`
	Application AppInfoApplication `xml:"application"`
}

type AppInfoApplication struct {
	XMLName xml.Name `xml:"application"`
	AppId   int      `xml:"app_id,attr"`
	AppName string   `xml:"app_name,attr"`
}

type BuildInfo struct {
	XMLName xml.Name       `xml:"buildinfo"`
	Build   BuildInfoBuild `xml:"build"`
}

type BuildInfoBuild struct {
	XMLName xml.Name `xml:"build"`
	Name    string   `xml:"version,attr"`
}

func (api API) populateReportDetailsFromAppInfo(appId int, detailedReport *DetailedReport) {
	var url = fmt.Sprintf("https://analysiscenter.veracode.com/api/5.0/getappinfo.do?app_id=%d", appId)
	response := api.makeApiRequest(url, http.MethodGet)

	data := AppInfo{}
	xml.Unmarshal(response, &data)

	detailedReport.AppId = data.Application.AppId
	detailedReport.AppName = data.Application.AppName
}

func (api API) populateReportDetailsFromBuildInfo(appId, buildId int, detailedReport *DetailedReport) {
	var url = fmt.Sprintf("https://analysiscenter.veracode.com/api/5.0/getbuildinfo.do?app_id=%d&build_id=%d", appId, buildId)
	response := api.makeApiRequest(url, http.MethodGet)

	data := BuildInfo{}
	xml.Unmarshal(response, &data)

	detailedReport.StaticAnalysis.ScanName = data.Build.Name
}
