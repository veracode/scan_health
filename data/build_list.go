package data

import (
	"encoding/xml"
	"fmt"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"html"
	"net/http"
	"sort"
)

type buildList struct {
	XMLName       xml.Name `xml:"buildlist"`
	AccountId     int      `xml:"account_id,attr"`
	ApplicationId int      `xml:"app_id,attr"`
	Builds        []build  `xml:"build"`
}

type build struct {
	BuildId  int    `xml:"build_id,attr"`
	ScanName string `xml:"version,attr"`
}

func (api API) populateBuildList(r *report.Report) {
	buildList := buildList{}

	sandboxFilter := ""

	// We may not have a sandbox id if no detailed report
	if r.Scan.SandboxId > 0 {
		sandboxFilter = fmt.Sprintf("&sandbox_id=%d", r.Scan.SandboxId)
	}

	var url = fmt.Sprintf("/api/5.0/getbuildlist.do?app_id=%d%s", r.Scan.ApplicationId, sandboxFilter)
	response := api.makeApiRequest(url, http.MethodGet)

	err := xml.Unmarshal(response, &buildList)

	if err != nil {
		utils.ErrorAndExit("Could not parse getbuildlist.do API response", err)
	}

	if buildList.ApplicationId != r.Scan.ApplicationId {
		utils.ErrorAndExit("Application Id mismatch", err)
	}

	var otherScans []report.Scan

	for _, b := range buildList.Builds {
		if r.Scan.BuildId == b.BuildId {
			if len(r.Scan.ScanName) == 0 {
				// Take the scan name from the build list because the detailed report may not yet be ready (scan not finished)
				r.Scan.ScanName = html.UnescapeString(b.ScanName)
			}
		} else {
			otherScans = append(otherScans, report.Scan{
				BuildId: b.BuildId,
			})
		}
	}

	// Sort by BuildId descending as there's nothing in the spec
	// to say what order they are in from getbuildlist
	sort.Slice(otherScans, func(i, j int) bool {
		return otherScans[i].BuildId > otherScans[j].BuildId
	})

	r.OtherScans = otherScans
}
