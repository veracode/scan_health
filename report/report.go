package report

import (
	"fmt"
	"time"

	"github.com/veracode/scan_health/v2/utils"
)

type healthTool struct {
	ReportDate time.Time `json:"report_date"`
	Version    string    `json:"version"`
	Region     string    `json:"region"`
}

type Scan struct {
	AccountId            int           `json:"account_id,omitempty"`
	BusinessUnit         string        `json:"business_unit,omitempty"`
	ApplicationId        int           `json:"application_id,omitempty"`
	ApplicationName      string        `json:"application_name,omitempty"`
	ScanName             string        `json:"scan_name,omitempty"`
	SandboxId            int           `json:"sandbox_id,omitempty"`
	SandboxName          string        `json:"sandbox_name,omitempty"`
	BuildId              int           `json:"build_id,omitempty"`
	AnalysisId           int           `json:"analysis_id,omitempty"`
	StaticAnalysisUnitId int           `json:"static_analysis_unit_id,omitempty"`
	ReviewModulesUrl     string        `json:"review_modules_url,omitempty"`
	TriageFlawsUrl       string        `json:"triage_flaws_url,omitempty"`
	TotalFilesUploaded   int           `json:"total_files_uploaded,omitempty"`
	TotalModules         int           `json:"total_modules,omitempty"`
	TotalModulesSelected int           `json:"total_modules_selected,omitempty"`
	EngineVersion        string        `json:"engine_version,omitempty"`
	SubmittedDate        time.Time     `json:"submitted_date,omitempty"`
	PublishedDate        time.Time     `json:"published_data,omitempty"`
	Duration             time.Duration `json:"duration,omitempty"`
	AnalysisSize         uint64        `json:"analysis_size,omitempty"`
	IsLatestScan         bool          `json:"is_latest_scan"`
	IsSCADataAvailable   bool          `json:"is_sca_data_available"`
}

type UploadedFile struct {
	Id           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Status       string `json:"status,omitempty"`
	MD5          string `json:"md5,omitempty"`
	IsThirdParty bool   `json:"is_third_party"`
	IsIgnored    bool   `json:"is_ignored"`
	Source       string `json:"source,omitempty"`
}

type issueSeverity string

const (
	IssueSeverityHigh   issueSeverity = "high"
	IssueSeverityMedium issueSeverity = "medium"
)

type Issue struct {
	Description     string        `json:"description,omitempty"`
	Severity        issueSeverity `json:"severity,omitempty"`
	AffectedFiles   []string      `json:"affected_files,omitempty"`
	AffectedModules []string      `json:"affected_modules,omitempty"`
}

type Report struct {
	HealthTool           healthTool     `json:"health_tool,omitempty"`
	LastAppActivity      time.Time      `json:"last_app_activity,omitempty"`
	LastSandboxActivity  time.Time      `json:"last_sandbox_activity,omitempty"`
	Scan                 Scan           `json:"scan,omitempty"`
	OtherScans           []Scan         `json:"-"`
	Flaws                FlawSummary    `json:"flaws,omitempty"`
	UploadedFiles        []UploadedFile `json:"uploaded_files,omitempty"`
	Modules              []Module       `json:"modules,omitempty"`
	SCAComponents        []string       `json:"sca_components,omitempty"`
	Issues               []Issue        `json:"issues,omitempty"`
	Recommendations      []string       `json:"recommendations,omitempty"`
	IsReportForOtherScan bool           `json:"-"`
}

func NewReport(buildId int, region, version string, isReportForOtherScan bool) *Report {
	return &Report{
		HealthTool:           healthTool{ReportDate: time.Now(), Version: version, Region: region},
		Scan:                 Scan{BuildId: buildId},
		Flaws:                FlawSummary{},
		UploadedFiles:        []UploadedFile{},
		Modules:              []Module{},
		Issues:               []Issue{},
		Recommendations:      []string{},
		IsReportForOtherScan: isReportForOtherScan,
	}
}

func (r *Report) ReportIssue(issue string, severity issueSeverity) {
	r.Issues = append(r.Issues, Issue{Description: issue, Severity: severity})
}

func (r *Report) GetReviewModulesUrl() string {
	return fmt.Sprintf("%s/auth/index.jsp#AnalyzeAppModuleList:%d:%d:%d:%d:%d::::%d",
		utils.ParseBaseUrlFromRegion(r.HealthTool.Region),
		r.Scan.AccountId,
		r.Scan.ApplicationId,
		r.Scan.BuildId,
		r.Scan.AnalysisId,
		r.Scan.StaticAnalysisUnitId,
		r.Scan.SandboxId)
}

func (r *Report) GetTriageFlawsUrl() string {
	return fmt.Sprintf("%s/auth/index.jsp#ReviewResultsStaticFlaws:%d:%d:%d:%d:%d::::%d",
		utils.ParseBaseUrlFromRegion(r.HealthTool.Region),
		r.Scan.AccountId,
		r.Scan.ApplicationId,
		r.Scan.BuildId,
		r.Scan.AnalysisId,
		r.Scan.StaticAnalysisUnitId,
		r.Scan.SandboxId)
}

func (r *Report) ReportFileIssue(issue string, severity issueSeverity, files []string) {
	for index, existingIssue := range r.Issues {
		if existingIssue.Description == issue {
			r.Issues[index].AffectedFiles = files
			return
		}
	}

	r.Issues = append(r.Issues, Issue{Description: issue, Severity: severity, AffectedFiles: files})
}

func (r *Report) ReportModuleIssue(issue string, severity issueSeverity, modules []string) {
	for index, existingIssue := range r.Issues {
		if existingIssue.Description == issue {
			r.Issues[index].AffectedModules = modules
			return
		}
	}

	r.Issues = append(r.Issues, Issue{Description: issue, Severity: severity, AffectedModules: modules})
}

func (r *Report) MakeRecommendation(recommendation string) {
	if !utils.IsStringInStringArray(recommendation, r.Recommendations) {
		r.Recommendations = append(r.Recommendations, recommendation)
	}
}

func (r *Report) Render(format, jsonFilePath string) {
	if jsonFilePath != "" {
		r.renderToJson(jsonFilePath)
	}

	if format == "json" {
		r.renderToJson("")
	} else {
		r.renderToConsole()
	}
}

func (r *Report) PrioritizeIssues() {
	var issues []Issue

	// Render the high severity issues first
	for _, issue := range r.Issues {
		if issue.Severity == IssueSeverityHigh {
			issues = append(issues, issue)
		}
	}

	for _, issue := range r.Issues {
		if issue.Severity != IssueSeverityHigh {
			issues = append(issues, issue)
		}
	}

	r.Issues = issues
}
