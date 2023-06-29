package report

import (
	"github.com/antfie/scan_health/v2/utils"
	"strings"
)

func (r *Report) GetSelectedModules() []Module {
	var selectedModules []Module

	for _, module := range r.Modules {
		if module.IsSelected {
			selectedModules = append(selectedModules, module)
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

	for _, module := range r.Modules {
		if module.IsSelected && utils.IsFileNameInFancyList(module.Name, fancyList) {
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
