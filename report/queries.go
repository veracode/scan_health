package report

import (
	"github.com/antfie/scan_health/v2/utils"
	"strings"
)

func (report *Report) FancyListMatchUploadedFiles(fancyList []string) []string {
	var foundFiles []string

	for _, uploadedFile := range report.UploadedFiles {

		if utils.IsFileNameInFancyList(uploadedFile.Name, fancyList) {
			if !utils.IsStringInStringArray(uploadedFile.Name, foundFiles) {
				foundFiles = append(foundFiles, uploadedFile.Name)
			}
		}
	}

	return foundFiles
}

func (report *Report) FancyListMatchModules(fancyList []string) []string {
	var selectedModules []string

	for _, module := range report.Modules {
		if utils.IsFileNameInFancyList(module.Name, fancyList) {
			if !utils.IsStringInStringArray(module.Name, selectedModules) {
				selectedModules = append(selectedModules, module.Name)
			}
		}
	}

	return selectedModules
}

func (report *Report) FancyListMatchSelectedModules(fancyList []string) []string {
	var selectedModules []string

	for _, module := range report.Modules {
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
