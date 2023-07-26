package report

import (
	"github.com/antfie/scan_health/v2/utils"
	"strings"
)

func (r *Report) GetSelectedModules() []Module {
	var selectedModules []Module

	for _, module := range r.Modules {
		for _, instance := range module.Instances {
			found := false

			for _, selectedModule := range selectedModules {
				if selectedModule.Name == module.Name {
					found = true
					continue
				}
			}

			if !found && instance.IsSelected {
				selectedModules = append(selectedModules, module)
			}
		}
	}

	return selectedModules
}

func (r *Report) FancyListMatchUploadedFiles(fancyList []string) []string {
	var foundFiles []string

	for _, uploadedFile := range r.UploadedFiles {

		if utils.IsFileNameInFancyList(uploadedFile.Name, fancyList) {
			if !utils.IsStringInStringArray(uploadedFile.Name, foundFiles) {
				foundFiles = append(foundFiles, uploadedFile.Name)
			}
		}
	}

	return foundFiles
}

func (r *Report) FancyListMatchModules(fancyList []string) []string {
	var selectedModules []string

	for _, module := range r.Modules {
		if utils.IsFileNameInFancyList(module.Name, fancyList) {
			if !utils.IsStringInStringArray(module.Name, selectedModules) {
				selectedModules = append(selectedModules, module.Name)
			}
		}
	}

	return selectedModules
}

func (r *Report) FancyListMatchSelectedModules(fancyList []string) []string {
	var selectedModules []string

	for _, module := range r.GetSelectedModules() {
		if utils.IsFileNameInFancyList(module.Name, fancyList) {
			if !utils.IsStringInStringArray(module.Name, selectedModules) {
				selectedModules = append(selectedModules, module.Name)
			}
		}
	}

	return selectedModules
}

func (module Module) IsDotNetOrCPPModule() bool {
	lowerCaseModuleName := strings.ToLower(module.Name)

	return strings.HasSuffix(lowerCaseModuleName, ".dll") || strings.HasSuffix(lowerCaseModuleName, ".exe")
}

func (module Module) IsJavaModule() bool {
	lowerCaseModuleName := strings.ToLower(module.Name)

	return strings.HasSuffix(lowerCaseModuleName, ".war") || strings.HasSuffix(lowerCaseModuleName, ".ear") || strings.HasSuffix(lowerCaseModuleName, ".jar")
}

func (module Module) IsJavaScriptModule() bool {
	return strings.HasPrefix(strings.ToLower(module.Name), "js files within")
}

func (module *Module) IsDependency() bool {
	for _, instance := range module.Instances {
		if instance.IsDependency {
			return true
		}
	}

	return false
}

func (module *Module) IsThirdParty() bool {
	for _, instance := range module.Instances {
		if instance.IsThirdParty {
			return true
		}
	}

	return false
}

func (module *Module) HasFatalErrors() bool {
	for _, instance := range module.Instances {
		if instance.HasFatalErrors {
			return true
		}
	}

	return false
}

func (module *Module) IsIgnored() bool {
	for _, instance := range module.Instances {
		if instance.IsIgnored {
			return true
		}
	}

	return false
}

func (module *Module) IsSelected() bool {
	for _, instance := range module.Instances {
		if instance.IsSelected {
			return true
		}
	}

	return false
}

func (module *Module) GetAllIssues() []string {
	var issues []string

	for _, instance := range module.Instances {
		for _, issue := range instance.Issues {
			for _, existingIssue := range issues {
				if strings.EqualFold(issue, existingIssue) {
					issues = append(issues, issue)
				}
			}
		}
	}

	return issues
}

func (module *Module) HasStatus(status string) bool {
	for _, instance := range module.Instances {
		if strings.EqualFold(status, instance.Status) {
			return true
		}
	}

	return false
}
