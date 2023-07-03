package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
	"strings"
)

func missingDebugSymbols(r *report.Report) {
	var foundModules []string

	for _, module := range r.Modules {
		// Only applicable for .net modules
		if !module.IsDotNetOrCPPModule() {
			continue
		}

		for _, issue := range module.Issues {
			if strings.Contains(issue, "No supporting files or PDB files") {
				if !utils.IsStringInStringArray(module.Name, foundModules) {
					foundModules = append(foundModules, module.Name)
				}
			}
		}
	}

	if len(foundModules) == 0 {
		return
	}

	message := fmt.Sprintf("A module was found to lack debug symbols (e.g., PDB files): %s.", foundModules[0])

	if len(foundModules) > 1 {
		message = fmt.Sprintf("%d modules were found lacking debug symbols (e.g., PDB files): %s.", len(foundModules), utils.Top5StringList(foundModules))
	}

	r.ReportIssue(message, report.IssueSeverityMedium)
	r.MakeRecommendation("Include PDB files for as many components as possible, especially first and second party components. This enables Veracode to accurately report line numbers for any flaws found within these components.")
}
