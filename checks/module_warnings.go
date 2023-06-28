package checks

import (
	"github.com/antfie/scan_health/v2/report"
)

func moduleWarnings(r *report.Report) {

}

//
//			if _, isMessageInMap := warnings[issue.Details]; !isMessageInMap {
//				warnings[issue.Details] = []string{}
//			}
//
//			if !utils.isStringInStringArray(module.Name, warnings[issue.Details]) {
//				warnings[issue.Details] = append(warnings[issue.Details], module.Name)
//			}
//		}
//
//		for _, statusMessage := range strings.Split(module.Status, ",") {
//			if module.Status == "OK" {
//				continue
//			}
//
//			formattedStatusMessage := strings.TrimSpace(statusMessage)
//
//			if strings.HasPrefix(formattedStatusMessage, "Unsupported Framework") {
//				// Ignore this
//				continue
//			}
//
//			if strings.HasPrefix(formattedStatusMessage, "Support Issue") {
//				// Ignore this
//				continue
//			}
//
//			if strings.HasPrefix(formattedStatusMessage, "PDB Files Missing") {
//				if strings.HasSuffix(formattedModuleName, ".dll") || strings.HasSuffix(formattedModuleName, ".exe") {
//					formattedStatusMessage = "Modules with missing PDB files"
//					data.makeRecommendation("Ensure you include PDB files for all first and second party .NET components. This enables Veracode to accurately report line numbers for any found flaws")
//				} else {
//					continue
//				}
//			}
//
//			if strings.Contains(formattedStatusMessage, "Missing Supporting Files") {
//				formattedStatusMessage = "Modules with missing supporting files"
//				data.makeRecommendation("Be sure to include all the components that make up the application within the upload. Do not omit any second or third-party libraries from the upload")
//			}
//
//			if !utils.isStringInStringArray(module.Name, warnings[formattedStatusMessage]) {
//				warnings[formattedStatusMessage] = append(warnings[formattedStatusMessage], module.Name)
//			}
//		}
//	}
//
