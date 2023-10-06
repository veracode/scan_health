package data

import (
	"encoding/xml"
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/fatih/color"
	"net/http"
	"sort"
	"time"
)

type build struct {
	BuildID           int       `xml:"build_id,attr"`
	PolicyUpdatedDate time.Time `xml:"policy_updated_date,attr"`
	Version           string    `xml:"version,attr"`
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

func findPreviousBuild(builds []build, targetBuildID int) *build {
	// Sort the builds by PolicyUpdatedDate as there's nothing in the spec
	// to say what order they are in from getbuildlist
	sort.Slice(builds, func(i, j int) bool {
		return builds[i].PolicyUpdatedDate.Before(builds[j].PolicyUpdatedDate)
	})

	var targetDate time.Time
	for _, build := range builds {
		if build.BuildID == targetBuildID {
			targetDate = build.PolicyUpdatedDate
			break
		}
	}

	if targetDate.IsZero() {
		return nil
	}

	var previousBuild *build
	for _, build := range builds {
		if build.PolicyUpdatedDate.Before(targetDate) {
			previousBuild = &build
		} else {
			break
		}
	}

	return previousBuild
}

func (api API) getPreviousBuildId(scan report.Scan) *buildList {
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

func (api API) populatePreviousBuildInfo(report *report.Report) {

	var buildList = api.getPreviousBuildId(report.Scan)

	if buildList == nil {
		color.HiYellow("Could not get list of builds, so cannot perform the check against the previous scan")
		return
	}
	var previousBuild = findPreviousBuild(buildList.Builds, report.Scan.BuildId)

	if previousBuild == nil {
		color.HiYellow("No previous build was found,  so cannot perform the check against the previous scan")
		return
	}

	// moduleList := api.retrievePrescanModuleListViaAPI(report.Scan.ApplicationId, previousBuild.BuildID)
}
