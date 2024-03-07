package report

import (
	"github.com/veracode/scan_health/v2/utils"
)

func (r *Report) GetSelectedModules() []Module {
	var selectedModules []Module

	for _, module := range r.Modules {
		found := false

		for _, instance := range module.Instances {
			if instance.Source == DetailedReportModuleSelected {
				found = true
			}
		}

		if found {
			selectedModules = append(selectedModules, module)
		}
	}

	return selectedModules
}

func (r *Report) GetPrescanModules() []Module {
	var prescanModules []Module

	for _, module := range r.Modules {
		for _, instance := range module.Instances {
			if !module.IsInListByName(prescanModules) && instance.Source == PrescanModuleList {
				prescanModules = append(prescanModules, module)
			}
		}
	}

	return prescanModules
}

func (r *Report) FancyListMatchUploadedFiles(fancyList []string) []string {
	var foundFiles []string

	for _, uploadedFile := range r.UploadedFiles {

		if utils.IsFileNameInFancyList(uploadedFile.Name, fancyList) {
			if !utils.IsStringInStringArray(uploadedFile.Name, foundFiles) && !uploadedFile.IsIgnored {
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
