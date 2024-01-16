package data

import (
	"encoding/xml"
	"fmt"
	"github.com/fatih/color"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"html"
	"net/http"
	"regexp"
	"strings"
)

type detailedReport struct {
	XMLName              xml.Name                     `xml:"detailedreport"`
	AccountId            int                          `xml:"account_id,attr"`
	AppId                int                          `xml:"app_id,attr"`
	AppName              string                       `xml:"app_name,attr"`
	SandboxId            int                          `xml:"sandbox_id,attr"`
	BuildId              int                          `xml:"build_id,attr"`
	AnalysisId           int                          `xml:"analysis_id,attr"`
	StaticAnalysisUnitId int                          `xml:"static_analysis_unit_id,attr"`
	TotalFlaws           int                          `xml:"total_flaws,attr"`
	UnmitigatedFlaws     int                          `xml:"flaws_not_mitigated,attr"`
	StaticAnalysis       detailedReportStaticAnalysis `xml:"static-analysis"`
	Flaws                []detailedReportFlaw         `xml:"severity>category>cwe>staticflaws>flaw"`
	IsLatestScan         bool                         `xml:"is_latest_build,attr"`
	BusinessUnit         string                       `xml:"business_unit,attr"`
	SCAResults           scaResults                   `xml:"software_composition_analysis"`
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
	ModulePath              string
}

type scaResults struct {
	XMLName              xml.Name                     `xml:"software_composition_analysis"`
	ServiceAvailable     string                       `xml:"sca_service_available,attr"` // sca_service_available = "false" when SCA is not enabled
	VulnerableComponents []detailedReportSCAComponent `xml:"vulnerable_components>component"`
}

type detailedReportSCAComponent struct {
	XMLName  xml.Name `xml:"component"`
	FileName string   `xml:"file_name,attr"`
}

func (api API) populateDetailedReport(r *report.Report) {
	detailedReport := detailedReport{}

	var url = fmt.Sprintf("/api/5.0/detailedreport.do?build_id=%d", r.Scan.BuildId)
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
	r.Scan.BusinessUnit = html.UnescapeString(detailedReport.BusinessUnit)
	r.Scan.ApplicationId = detailedReport.AppId
	r.Scan.ApplicationName = html.UnescapeString(detailedReport.AppName)
	r.Scan.SandboxId = detailedReport.SandboxId
	r.Scan.ScanName = html.UnescapeString(detailedReport.StaticAnalysis.ScanName)
	r.Scan.ReviewModulesUrl = detailedReport.getReviewModulesUrl(r.HealthTool.Region)
	r.Scan.TriageFlawsUrl = detailedReport.getTriageFlawsUrl(r.HealthTool.Region)
	r.Scan.EngineVersion = detailedReport.StaticAnalysis.EngineVersion
	r.Scan.AnalysisSize = detailedReport.StaticAnalysis.AnalysisSizeBytes

	if len(detailedReport.StaticAnalysis.SubmittedDate) > 0 {
		r.Scan.SubmittedDate = utils.ParseVeracodeDate(detailedReport.StaticAnalysis.SubmittedDate).Local()
	}

	if len(detailedReport.StaticAnalysis.PublishedDate) > 0 {
		r.Scan.PublishedDate = utils.ParseVeracodeDate(detailedReport.StaticAnalysis.PublishedDate).Local()
	}

	if detailedReport.StaticAnalysis.SubmittedDate != "" && detailedReport.StaticAnalysis.PublishedDate != "" {
		r.Scan.ScanDuration = r.Scan.PublishedDate.Sub(r.Scan.SubmittedDate)
	}

	r.Scan.IsLatestScan = detailedReport.IsLatestScan

	populateDetailedReportModules(r, detailedReport.StaticAnalysis)
	populateModulesFromFlaws(r, detailedReport)
	populateFlawSummaries(r, detailedReport)
	populateThirdPartyFiles(r, detailedReport)
}

func populateDetailedReportModules(r *report.Report, staticAnalysis detailedReportStaticAnalysis) {
	for _, module := range staticAnalysis.Modules {
		r.AddModuleInstance(normalizeModuleName(html.UnescapeString(module.Name)), report.ModuleInstance{
			Compiler:        html.UnescapeString(module.Compiler),
			OperatingSystem: html.UnescapeString(module.Os),
			Architecture:    html.UnescapeString(module.Architecture),
			IsSelected:      true,
			WasScanned:      true,
			Source:          "detailed_report_module_selected",
		})
	}
}

func normalizeModuleName(moduleName string) string {
	// Sometimes the detailed report will return [modulename]_htmljscode.veracodegen.htmla.jsa
	// rather than JS files within [modulename]. We should use the latter, to align with the
	// prescan data.
	pattern := `(.+)_htmljscode\.veracodegen\.htmla\.jsa`
	re := regexp.MustCompile(pattern)
	if matches := re.FindStringSubmatch(moduleName); matches != nil {
		return fmt.Sprintf("JS files within %s", matches[1])
	}

	return moduleName
}

func populateModulesFromFlaws(r *report.Report, detailedReport detailedReport) {
	for index, flaw := range detailedReport.Flaws {
		detailedReport.Flaws[index].Module = html.UnescapeString(flaw.Module)

		// Set the module path e.g. /a.war/b.jar/c
		detailedReport.Flaws[index].ModulePath = flaw.Module

		if strings.Contains(flaw.Module, "/") {
			modulePathParts := strings.Split(flaw.Module, "/")
			detailedReport.Flaws[index].Module = modulePathParts[len(modulePathParts)-1]

			// Also update the local copy of this flaw
			flaw.Module = detailedReport.Flaws[index].Module

			for modulePartIndex, modulePart := range modulePathParts {
				r.AddModuleInstance(modulePart, report.ModuleInstance{
					WasScanned:   true,
					IsDependency: modulePartIndex > 0,
					Source:       "detailed_report_module_derived_from_flaw_module_path",
				})

				if modulePartIndex > 0 {
					r.AddModuleDependency(modulePart, modulePathParts[modulePartIndex-1])
				}
			}
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
		// Update report totals
		r.Flaws.Total++

		if flaw.AffectsPolicyCompliance {
			r.Flaws.TotalAffectingPolicy++
		}

		if flaw.isOpen() {
			if flaw.AffectsPolicyCompliance {
				r.Flaws.OpenAffectingPolicy++
			} else {
				r.Flaws.OpenButNotAffectingPolicy++
			}
		} else if flaw.isMitigated() {
			r.Flaws.Mitigated++
		} else if flaw.isFixed() {
			r.Flaws.Fixed++
		}

		// Update totals per-module affected
		for moduleIndex, module := range r.Modules {
			// For each module in the module path
			modulePathParts := strings.Split(flaw.ModulePath, "/")
			for _, modulePath := range modulePathParts {
				if strings.EqualFold(modulePath, module.Name) {
					flawSummary := &r.Modules[moduleIndex].Flaws
					flawSummary.Total++

					if flaw.AffectsPolicyCompliance {
						flawSummary.TotalAffectingPolicy++
					}

					if flaw.isOpen() {
						if flaw.AffectsPolicyCompliance {
							flawSummary.OpenAffectingPolicy++
						} else {
							flawSummary.OpenButNotAffectingPolicy++
						}
					} else if flaw.isMitigated() {
						flawSummary.Mitigated++
					} else if flaw.isFixed() {
						flawSummary.Fixed++
					}
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

func populateThirdPartyFiles(r *report.Report, detailedReport detailedReport) {
	r.Scan.IsSCADataAvailable = strings.ToLower(detailedReport.SCAResults.ServiceAvailable) != "false"

	for _, component := range detailedReport.SCAResults.VulnerableComponents {
		formattedComponentName := html.UnescapeString(component.FileName)

		if !utils.IsStringInStringArray(formattedComponentName, r.SCAComponents) {
			r.SCAComponents = append(r.SCAComponents, formattedComponentName)
		}
	}
}
