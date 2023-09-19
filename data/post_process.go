package data

import (
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
)

func postProcessData(r *report.Report) {
	postProcessThirdParty(r)
}

func postProcessThirdParty(r *report.Report) {
	var thirdPartyFiles []string

	for _, uploadedFile := range r.UploadedFiles {
		if uploadedFile.Source == "detailed_report_sca_component" && !utils.IsStringInStringArray(uploadedFile.Name, thirdPartyFiles) {
			thirdPartyFiles = append(thirdPartyFiles, uploadedFile.Name)
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
