package report

import (
	"github.com/veracode/scan_health/v2/utils"
	"strings"
)

type Module struct {
	Name         string           `json:"name,omitempty"`
	IsThirdParty bool             `json:"is_third_party"`
	IsIgnored    bool             `json:"is_ignored"`
	Instances    []ModuleInstance `json:"instances,omitempty"`
	Flaws        FlawSummary      `json:"flaws,omitempty"`
	FlawDetails  []FlawDetails    `json:"-"`
	DependencyOf []string         `json:"dependency_of,omitempty"`
}

type ModuleInstanceSource string

const (
	DetailedReportModuleSelected ModuleInstanceSource = "detailed_report_module_selected"
	PrescanModuleList            ModuleInstanceSource = "prescan_module_list"
)

type ModuleInstance struct {
	Id              int                  `json:"id,omitempty"`
	Compiler        string               `json:"compiler,omitempty"`
	OperatingSystem string               `json:"operating_system,omitempty"`
	Architecture    string               `json:"architecture,omitempty"`
	IsDependency    bool                 `json:"is_dependency"`
	IsSelected      bool                 `json:"is_selected"`
	WasScanned      bool                 `json:"was_scanned"`
	HasFatalErrors  bool                 `json:"has_fatal_errors"`
	Status          string               `json:"status,omitempty"`
	Platform        string               `json:"platform,omitempty"`
	Size            string               `json:"size,omitempty"`
	MD5             string               `json:"md5,omitempty"`
	Issues          []string             `json:"issues,omitempty"`
	SizeBytes       int                  `json:"size_bytes,omitempty"`
	Source          ModuleInstanceSource `json:"source,omitempty"`
}

func (m *Module) IsSelected() bool {
	for _, i := range m.Instances {
		if i.IsSelected {
			return true
		}
	}

	return false
}

func (m *Module) IsDotNetOrCPPModule() bool {
	lowerCaseModuleName := strings.ToLower(m.Name)

	return strings.HasSuffix(lowerCaseModuleName, ".dll") || strings.HasSuffix(lowerCaseModuleName, ".exe")
}

func (m *Module) IsJavaModule() bool {
	lowerCaseModuleName := strings.ToLower(m.Name)

	return strings.HasSuffix(lowerCaseModuleName, ".war") || strings.HasSuffix(lowerCaseModuleName, ".ear") || strings.HasSuffix(lowerCaseModuleName, ".jar")
}

func (m *Module) IsJavaScriptModule() bool {
	return strings.HasPrefix(strings.ToLower(m.Name), "js files within") || strings.HasPrefix(strings.ToLower(m.Name), "js files extracted from")
}

func (m *Module) IsNodeModule() bool {
	return strings.Contains(strings.ToLower(m.Name), "_nodemodule_")
}

func (m *Module) IsDependency() bool {
	for _, instance := range m.Instances {
		if instance.IsDependency {
			return true
		}
	}

	return false
}

func (m *Module) HasFatalErrors() bool {
	for _, instance := range m.Instances {
		if instance.HasFatalErrors {
			return true
		}
	}

	return false
}

func (m *Module) Issues() []string {
	var issues []string

	for _, instance := range m.Instances {
		for _, issue := range instance.Issues {
			if !utils.IsStringInStringArray(issue, issues) {
				issues = append(issues, issue)
			}
		}
	}

	return issues
}

func (m *Module) HasStatus(status string) bool {
	for _, instance := range m.Instances {
		if strings.EqualFold(status, instance.Status) {
			return true
		}
	}

	return false
}

func (m *Module) IsInListByName(modules []Module) bool {
	for _, existingModule := range modules {
		if existingModule.Name == m.Name {
			return true
		}
	}

	return false
}
