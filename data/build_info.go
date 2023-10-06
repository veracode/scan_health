package data

import (
	"encoding/xml"
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
	"net/http"
)

type buildInfo struct {
	XMLName xml.Name       `xml:"buildinfo"`
	Build   buildInfoBuild `xml:"build"`
}

type buildInfoBuild struct {
	XMLName xml.Name `xml:"build"`
	Version string   `xml:"version,attr"`
}

func (api API) populateBuildInfo(report *report.Report) {
	var url = fmt.Sprintf("https://analysiscenter.veracode.com/api/5.0/getbuildinfo.do?app_id=%d&build_id=%d", report.Scan.ApplicationId, report.Scan.BuildId)
	response := api.makeApiRequest(url, http.MethodGet)

	data := buildInfo{}
	err := xml.Unmarshal(response, &data)

	if err != nil {
		utils.ErrorAndExit("Could not parse response from getbuildinfo.do API response", err)
	}

	url = fmt.Sprintf("https://analysiscenter.veracode.com/api/5.0/getprescanresults.do?app_id=%d&build_id=%d", report.Scan.ApplicationId, report.Scan.BuildId)
	response = api.makeApiRequest(url, http.MethodGet)

	moduleList := prescanModuleList{}

	err = xml.Unmarshal(response, &moduleList)

	if err != nil {
		utils.ErrorAndExit("Could not get prescan results", err)
	}

	populateModuleInstances(report, moduleList)
}
