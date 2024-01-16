package data

import (
	"encoding/xml"
	"fmt"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"net/http"
)

type sandboxList struct {
	XMLName   xml.Name      `xml:"sandboxlist"`
	Sandboxes []sandboxInfo `xml:"sandbox"`
}

type sandboxInfo struct {
	XMLName      xml.Name `xml:"sandbox"`
	Id           int      `xml:"sandbox_id,attr"`
	Name         string   `xml:"sandbox_name,attr"`
	ModifiedDate string   `xml:"last_modified,attr"`
}

func (api API) populateSandboxInfo(report *report.Report) {
	var url = fmt.Sprintf("/api/5.0/getsandboxlist.do?app_id=%d", report.Scan.ApplicationId)
	response := api.makeApiRequest(url, http.MethodGet)

	data := sandboxList{}
	err := xml.Unmarshal(response, &data)

	if err != nil {
		utils.ErrorAndExit("Could not parse response from getsandboxlist.do API response", err)
	}

	for _, sandbox := range data.Sandboxes {
		if sandbox.Id == report.Scan.SandboxId {
			report.Scan.SandboxName = sandbox.Name

			if len(sandbox.ModifiedDate) > 0 {
				report.LastSandboxActivity = utils.ParseVeracodeDate(sandbox.ModifiedDate).Local()
			}
		}
	}
}
