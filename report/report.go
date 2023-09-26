package report

import (
	"time"

	"github.com/antfie/scan_health/v2/utils"
)

type HealthTool struct {
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
	ReviewModulesUrl     string        `json:"review_modules_url,omitempty"`
	TriageFlawsUrl       string        `json:"triage_flaws_url,omitempty"`
	TotalFilesUploaded   int           `json:"total_files_uploaded,omitempty"`
	TotalModules         int           `json:"total_modules,omitempty"`
	TotalModulesSelected int           `json:"total_modules_selected,omitempty"`
	EngineVersion        string        `json:"engine_version,omitempty"`
	SubmittedDate        time.Time     `json:"submitted_date,omitempty"`
	PublishedDate        time.Time     `json:"published_data,omitempty"`
	ScanDuration         time.Duration `json:"scan_duration,omitempty"`
	AnalysisSize         uint64        `json:"analysis_size,omitempty"`
	IsLatestScan         bool          `json:"is_latest_scan"`
	IsSCADataAvailable   bool          `json:"is_sca_data_available"`
}

type FlawSummary struct {
	Total                     int `json:"total,omitempty"`
	Fixed                     int `json:"fixed,omitempty"`
	TotalAffectingPolicy      int `json:"total_affecting_policy,omitempty"`
	Mitigated                 int `json:"mitigated,omitempty"`
	OpenAffectingPolicy       int `json:"open_affecting_policy,omitempty"`
	OpenButNotAffectingPolicy int `json:"open_but_not_affecting_policy,omitempty"`
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

type Module struct {
	Name         string           `json:"name,omitempty"`
	IsThirdParty bool             `json:"is_third_party"`
	IsIgnored    bool             `json:"is_ignored"`
	Instances    []ModuleInstance `json:"instances,omitempty"`
	Flaws        FlawSummary      `json:"flaws,omitempty"`
	DependencyOf []string         `json:"dependency_of,omitempty"`
}

type ModuleInstance struct {
	Id              int      `json:"id,omitempty"`
	Compiler        string   `json:"compiler,omitempty"`
	OperatingSystem string   `json:"operating_system,omitempty"`
	Architecture    string   `json:"architecture,omitempty"`
	IsDependency    bool     `json:"is_dependency"`
	IsSelected      bool     `json:"is_selected"`
	WasScanned      bool     `json:"was_scanned"`
	HasFatalErrors  bool     `json:"has_fatal_errors"`
	Status          string   `json:"status,omitempty"`
	Platform        string   `json:"platform,omitempty"`
	Size            string   `json:"size,omitempty"`
	MD5             string   `json:"md5,omitempty"`
	Issues          []string `json:"issues,omitempty"`
	SizeBytes       int      `json:"size_bytes,omitempty"`
	Source          string   `json:"source,omitempty"`
}

type IssueSeverity string

const (
	IssueSeverityHigh   IssueSeverity = "high"
	IssueSeverityMedium IssueSeverity = "medium"
)

type Issue struct {
	Description     string        `json:"description,omitempty"`
	Severity        IssueSeverity `json:"severity,omitempty"`
	AffectedFiles   []string      `json:"affected_files,omitempty"`
	AffectedModules []string      `json:"affected_modules,omitempty"`
}

type Report struct {
	HealthTool          HealthTool     `json:"health_tool,omitempty"`
	LastAppActivity     time.Time      `json:"last_app_activity,omitempty"`
	LastSandboxActivity time.Time      `json:"last_sandbox_activity,omitempty"`
	Scan                Scan           `json:"scan,omitempty"`
	Flaws               FlawSummary    `json:"flaws,omitempty"`
	UploadedFiles       []UploadedFile `json:"uploaded_files,omitempty"`
	Modules             []Module       `json:"modules,omitempty"`
	SCAComponents       []string       `json:"sca_components,omitempty"`
	Issues              []Issue        `json:"issues,omitempty"`
	Recommendations     []string       `json:"recommendations,omitempty"`
}

func NewReport(buildId int, region, version string) *Report {
	return &Report{
		HealthTool:      HealthTool{ReportDate: time.Now(), Version: version, Region: region},
		Scan:            Scan{BuildId: buildId},
		Flaws:           FlawSummary{},
		UploadedFiles:   []UploadedFile{},
		Modules:         []Module{},
		Issues:          []Issue{},
		Recommendations: []string{},
	}
}

func (r *Report) ReportIssue(issue string, severity IssueSeverity) {
	r.Issues = append(r.Issues, Issue{Description: issue, Severity: severity})
}

func (r *Report) ReportFileIssue(issue string, severity IssueSeverity, files []string) {
	for index, existingIssue := range r.Issues {
		if existingIssue.Description == issue {
			r.Issues[index].AffectedFiles = files
			return
		}
	}

	r.Issues = append(r.Issues, Issue{Description: issue, Severity: severity, AffectedFiles: files})
}

func (r *Report) ReportModuleIssue(issue string, severity IssueSeverity, modules []string) {
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
