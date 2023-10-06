package data

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"net/http"
	"sort"
	"time"
)

type build struct {
	BuildID           int       `xml:"build_id,attr"`
	PolicyUpdatedDate time.Time `xml:"policy_updated_date,attr"`
	Version           string    `xml:"version,attr"`
	DynamicScanType   string    `xml:"dynamic_scan_type,attr"`
}

type buildList struct {
	XMLName          xml.Name `xml:"buildlist"`
	AccountID        int      `xml:"account_id,attr"`
	AppID            int      `xml:"app_id,attr"`
	AppName          string   `xml:"app_name,attr"`
	BuildListVersion string   `xml:"buildlist_version,attr"`
	SchemaLocation   string   `xml:"schemaLocation,attr"`
	Builds           []build  `xml:"build"`
}

func findPreviousBuildIdInBuildList(builds []build, currentBuildID int) (int, error) {
	// Find current build ID
	var targetDate time.Time
	for _, build := range builds {
		if build.BuildID == currentBuildID {
			targetDate = build.PolicyUpdatedDate
			break
		}
	}

	if targetDate.IsZero() {
		return 0, errors.New("could not find Build ID")
	}

	// Sort the builds by PolicyUpdatedDate as there's nothing in the spec
	// to say what order they are in from getbuildlist
	sort.Slice(builds, func(i, j int) bool {
		return builds[i].PolicyUpdatedDate.Before(builds[j].PolicyUpdatedDate)
	})

	var previousBuildId int
	for _, build := range builds {

		if build.DynamicScanType != "" {
			continue
		}

		if build.PolicyUpdatedDate.Before(targetDate) {
			previousBuildId = build.BuildID
		} else {
			break
		}
	}

	return previousBuildId, nil
}

func (api API) getBuildList(scan report.Scan) *buildList {
	var sandboxId = scan.SandboxId
	var appId = scan.ApplicationId

	buildList := buildList{}

	var url = fmt.Sprintf("https://analysiscenter.veracode.com/api/5.0/getbuildlist.do?app_id=%d&sandbox_id=%d", appId, sandboxId)
	response := api.makeApiRequest(url, http.MethodGet)

	err := xml.Unmarshal(response, &buildList)

	if err != nil {
		return nil
	}

	return &buildList
}

func (api API) GetPreviousBuildId(report *report.Report) (int, error) {

	var buildList = api.getBuildList(report.Scan)

	if buildList == nil {
		return 0, errors.New(`could not get list of builds, so cannot perform the check against the previous scan`)
	}

	previousBuildId, err := findPreviousBuildIdInBuildList(buildList.Builds, report.Scan.BuildId)

	if err != nil {
		return 0, errors.New(`No previous build was found,  so cannot perform the check against the previous scan`)
	}

	return previousBuildId, nil
}
