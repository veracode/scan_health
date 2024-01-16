package data

import (
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"sort"
)

func postProcessData(r *report.Report) {
	postProcessThirdParty(r)
	sortData(r)
}

func postProcessThirdParty(r *report.Report) {
	var thirdPartyFiles []string

	for _, component := range r.SCAComponents {
		if !utils.IsStringInStringArray(component, thirdPartyFiles) {
			thirdPartyFiles = append(thirdPartyFiles, component)
		}
	}

	for index, uploadedFile := range r.UploadedFiles {
		if !uploadedFile.IsThirdParty && utils.IsStringInStringArray(uploadedFile.Name, thirdPartyFiles) {
			r.UploadedFiles[index].IsThirdParty = true
		}
	}

	for index, module := range r.Modules {
		if !module.IsThirdParty && utils.IsStringInStringArray(module.Name, thirdPartyFiles) {
			r.Modules[index].IsThirdParty = true
		}
	}
}

func sortData(r *report.Report) {
	// Sort uploaded files by name for consistency
	sort.Slice(r.UploadedFiles, func(i, j int) bool {
		return r.UploadedFiles[i].Name < r.UploadedFiles[j].Name
	})

	// Sort modules by name for consistency
	sort.Slice(r.Modules, func(i, j int) bool {
		return r.Modules[i].Name < r.Modules[j].Name
	})
}
