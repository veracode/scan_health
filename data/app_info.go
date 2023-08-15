package data

import (
	"encoding/xml"
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
	"net/http"
)

type appInfo struct {
	XMLName     xml.Name           `xml:"appinfo"`
	Application appInfoApplication `xml:"application"`
}

type appInfoApplication struct {
	XMLName xml.Name `xml:"application"`
	AppId   int      `xml:"app_id,attr"`
	AppName string   `xml:"app_name,attr"`
	Updated string   `xml:"policy_updated_date,attr"`
}

func (api API) populateAppInfo(report *report.Report) {
	url := fmt.Sprintf("https://analysiscenter.veracode.com/api/5.0/getappinfo.do?app_id=%d", report.Scan.ApplicationId)
	response := api.makeApiRequest(url, http.MethodGet)

	data := appInfo{}
	err := xml.Unmarshal(response, &data)

	if err != nil {
		utils.ErrorAndExit("Could not parse getappinfo.do API response", err)
	}

	report.Scan.ApplicationId = data.Application.AppId
	report.Scan.ApplicationName = data.Application.AppName

	if len(data.Application.Updated) > 0 {
		report.LastAppActivity = utils.ParseVeracodeDate(data.Application.Updated).Local()
	}
}
