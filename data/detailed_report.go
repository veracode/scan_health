package data

import (
	"encoding/xml"
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
	"github.com/fatih/color"
	"net/http"
	"strings"
)

type detailedReport struct {
	XMLName              xml.Name                     `xml:"detailedreport"`
	AccountId            int                          `xml:"account_id,attr"`
	AppId                int                          `xml:"app_id,attr"`
	AppName              string                       `xml:"app_name,attr"`
	SandboxId            int                          `xml:"sandbox_id,attr"`
	SandboxName          string                       `xml:"sandbox_name,attr"`
	BuildId              int                          `xml:"build_id,attr"`
	AnalysisId           int                          `xml:"analysis_id,attr"`
	StaticAnalysisUnitId int                          `xml:"static_analysis_unit_id,attr"`
	TotalFlaws           int                          `xml:"total_flaws,attr"`
	UnmitigatedFlaws     int                          `xml:"flaws_not_mitigated,attr"`
	StaticAnalysis       detailedReportStaticAnalysis `xml:"static-analysis"`
	Flaws                []detailedReportFlaw         `xml:"severity>category>cwe>staticflaws>flaw"`
	IsLatestScan         bool                         `xml:"is_latest_build,attr"`
	BusinessUnit         string                       `xml:"business_unit,attr"`
}

type detailedReportStaticAnalysis struct {
	XMLName           xml.Name               `xml:"static-analysis"`
	EngineVersion     string                 `xml:"engine_version,attr"`
	SubmittedDate     string                 `xml:"submitted_date,attr"`
	PublishedDate     string                 `xml:"published_date,attr"`
	ScanName          string                 `xml:"version,attr"`
	Score             int                    `xml:"score,attr"`
	AnalysisSizeBytes uint64                 `xml:"analysis_size_bytes,attr"`
	Modules           []detailedReportModule `xml:"modules>module"`
}

type detailedReportModule struct {
	XMLName      xml.Name `xml:"module"`
	Name         string   `xml:"name,attr"`
	Compiler     string   `xml:"compiler,attr"`
	Os           string   `xml:"os,attr"`
	Architecture string   `xml:"architecture,attr"`
	IsIgnored    bool
	IsThirdParty bool
}

type detailedReportFlaw struct {
	XMLName                 xml.Name `xml:"flaw"`
	ID                      int      `xml:"issueid,attr"`
	CWE                     int      `xml:"cweid,attr"`
	AffectsPolicyCompliance bool     `xml:"affects_policy_compliance,attr"`
	Module                  string   `xml:"module,attr"`
	RemediationStatus       string   `xml:"remediation_status,attr"`     // Fixed, New, Reopened, Mitigated, Potential False Positive
	MitigationStatus        string   `xml:"mitigation_status,attr"`      // none, accepted, rejected
	Mitigation              string   `xml:"mitigation_status_desc,attr"` // Mitigation Accepted, Not Mitigated, Mitigation Proposed
}

func (api API) populateDetailedReport(r *report.Report) {
	detailedReport := detailedReport{}

	var url = fmt.Sprintf("https://analysiscenter.veracode.com/api/5.0/detailedreport.do?build_id=%d", r.Scan.BuildId)
	response := api.makeApiRequest(url, http.MethodGet)

	if strings.Contains(string(response[:]), "<error>No report available.</error>") {
		color.HiYellow(fmt.Sprintf("Warning: There was no detailed report for Build id %d. Has the scan finished?", r.Scan.BuildId))
		return
	}

	err := xml.Unmarshal(response, &detailedReport)

	if err != nil {
		utils.ErrorAndExit("Could not parse detailedreport.do API response", err)
	}

	r.Scan.AccountId = detailedReport.AccountId
	r.Scan.BusinessUnit = detailedReport.BusinessUnit
	r.Scan.ApplicationId = detailedReport.AppId
	r.Scan.ApplicationName = detailedReport.AppName
	r.Scan.SandboxId = detailedReport.SandboxId
	r.Scan.SandboxName = detailedReport.SandboxName
	r.Scan.ScanName = detailedReport.StaticAnalysis.ScanName
	r.Scan.ReviewModulesUrl = detailedReport.getReviewModulesUrl(r.HealthTool.Region)
	r.Scan.TriageFlawsUrl = detailedReport.getTriageFlawsUrl(r.HealthTool.Region)
	r.Scan.EngineVersion = detailedReport.StaticAnalysis.EngineVersion
	r.Scan.AnalysisSize = detailedReport.StaticAnalysis.AnalysisSizeBytes
	r.Scan.SubmittedDate = utils.ParseVeracodeDate(detailedReport.StaticAnalysis.SubmittedDate).Local()
	r.Scan.PublishedDate = utils.ParseVeracodeDate(detailedReport.StaticAnalysis.PublishedDate).Local()
	r.Scan.ScanDuration = r.Scan.PublishedDate.Sub(r.Scan.SubmittedDate)
	r.Scan.IsLatestScan = detailedReport.IsLatestScan

	if r.Scan.IsLatestScan {
		r.LastAppActivity = r.Scan.PublishedDate
	} else {
		color.HiYellow("Warning: This is not the latest scan")
	}

	populateDetailedReportModules(r, detailedReport.StaticAnalysis)
	populateModulesFromFlaws(r, detailedReport)
	populateFlawSummaries(r, detailedReport)
}

func populateDetailedReportModules(r *report.Report, staticAnalysis detailedReportStaticAnalysis) {
	for _, module := range staticAnalysis.Modules {
		r.Modules = append(r.Modules, report.Module{
			Name:         module.Name,
			Compiler:     module.Compiler,
			Os:           module.Os,
			Architecture: module.Architecture,
			IsSelected:   true,
		})
	}
}

func populateModulesFromFlaws(r *report.Report, detailedReport detailedReport) {
	for index, flaw := range detailedReport.Flaws {
		isDependentModule := false

		if strings.Contains(flaw.Module, "/") {
			modulePathParts := strings.Split(flaw.Module, "/")
			detailedReport.Flaws[index].Module = modulePathParts[len(modulePathParts)-1]

			// Also update the local copy of this flaw
			flaw.Module = detailedReport.Flaws[index].Module

			isDependentModule = true
		}

		moduleFound := false

		for _, module := range r.Modules {
			if strings.EqualFold(flaw.Module, module.Name) {
				moduleFound = true
				break
			}
		}

		if !moduleFound {
			r.Modules = append(r.Modules, report.Module{
				Name:         flaw.Module,
				IsDependency: isDependentModule,
			})
		}
	}
}

func (report detailedReport) getReviewModulesUrl(region string) string {
	return fmt.Sprintf("%s/auth/index.jsp#AnalyzeAppModuleList:%d:%d:%d:%d:%d::::%d",
		utils.ParseBaseUrlFromRegion(region),
		report.AccountId,
		report.AppId,
		report.BuildId,
		report.AnalysisId,
		report.StaticAnalysisUnitId,
		report.SandboxId)
}

func (report detailedReport) getTriageFlawsUrl(region string) string {
	return fmt.Sprintf("%s/auth/index.jsp#ReviewResultsStaticFlaws:%d:%d:%d:%d:%d::::%d",
		utils.ParseBaseUrlFromRegion(region),
		report.AccountId,
		report.AppId,
		report.BuildId,
		report.AnalysisId,
		report.StaticAnalysisUnitId,
		report.SandboxId)
}

func populateFlawSummaries(r *report.Report, detailedReport detailedReport) {
	for _, flaw := range detailedReport.Flaws {
		for _, module := range r.Modules {
			if strings.EqualFold(flaw.Module, module.Name) {
				module.Flaws.Total++
				r.Flaws.Total++

				if flaw.AffectsPolicyCompliance {
					module.Flaws.TotalAffectingPolicy++
					r.Flaws.TotalAffectingPolicy++
				}

				if flaw.isOpen() {
					if flaw.AffectsPolicyCompliance {
						module.Flaws.OpenAffectingPolicy++
						r.Flaws.OpenAffectingPolicy++
					} else {
						module.Flaws.OpenButNotAffectingPolicy++
						r.Flaws.OpenButNotAffectingPolicy++
					}
				} else if flaw.isMitigated() {
					module.Flaws.Mitigated++
					r.Flaws.Mitigated++
				} else if flaw.isFixed() {
					module.Flaws.Fixed++
					r.Flaws.Fixed++
				}
			}
		}
	}
}

func (report detailedReport) getPolicyAffectingFlawCount() int {
	var count = 0

	for _, flaw := range report.Flaws {
		if flaw.AffectsPolicyCompliance {
			count++
		}
	}

	return count
}

func (flaw detailedReportFlaw) isOpen() bool {
	return !flaw.isFixed() && !flaw.isMitigated()
}

func (flaw detailedReportFlaw) isFixed() bool {
	return flaw.RemediationStatus == "Fixed"
}

func (flaw detailedReportFlaw) isMitigated() bool {
	return !(flaw.MitigationStatus == "none" || flaw.MitigationStatus == "rejected")
}
