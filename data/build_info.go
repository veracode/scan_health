package data

import (
	"encoding/xml"
	"fmt"
	"github.com/antfie/scan_health/v2/utils"
	"net/http"
)

type buildInfo struct {
	XMLName       xml.Name       `xml:"buildinfo"`
	ApplicationId int            `xml:"app_id,attr"`
	Build         buildInfoBuild `xml:"build"`
}

type buildInfoBuild struct {
	XMLName      xml.Name     `xml:"build"`
	AnalysisUnit analysisUnit `xml:"analysis_unit"`
}

type analysisUnit struct {
	XMLName       xml.Name `xml:"analysis_unit"`
	Status        string   `xml:"status,attr"`
	PublishedDate string   `xml:"published_date,attr"`
}

func (api API) GetBuildInfo(applicationId, buildId int) *buildInfo {
	var url = fmt.Sprintf("/api/5.0/getbuildinfo.do?app_id=%d&build_id=%d", applicationId, buildId)
	response := api.makeApiRequest(url, http.MethodGet)

	buildInfo := &buildInfo{}
	err := xml.Unmarshal(response, &buildInfo)

	if err != nil {
		utils.ErrorAndExit("Could not parse response from getbuildinfo.do API response", err)
	}

	if buildInfo.ApplicationId != applicationId {
		utils.ErrorAndExit("Application Id mismatch", err)
	}

	return buildInfo
}
