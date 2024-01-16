package checks

import (
	"fmt"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"strings"
)

func releaseBuild(r *report.Report) {

	// TODO: we cannot do this test until we get file path data from the API
	//detectDotNetReleasePathsInModuleIssues(r)
}

// Test cases:
// https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:71832:1712306:26134566:26105587:26121237::::5355525

func detectDotNetReleasePathsInModuleIssues(r *report.Report) {
	var foundModules []string

	for _, module := range r.Modules {
		// Only applicable for .net modules
		if !module.IsDotNetOrCPPModule() {
			continue
		}

		// Ignore junk
		if module.IsIgnored || module.IsThirdParty {
			continue
		}

		for _, issue := range module.GetAllIssues() {
			if strings.Contains(strings.ToLower(issue), "release/") {
				if !utils.IsStringInStringArray(module.Name, foundModules) {
					foundModules = append(foundModules, module.Name)
				}
			}
		}
	}

	if len(foundModules) == 0 {
		return
	}

	message := fmt.Sprintf("A module appeared to contain or depend on components that were compiled for release: \"%s\". Veracode expects (where possible) a debug build for optimal scan quality and accurate line number reporting.", foundModules[0])

	if len(foundModules) > 1 {
		message = fmt.Sprintf(
			"%d modules appeared to contain or depend on on components that were compiled for release: %s. Veracode expects (where possible) a debug build for optimal scan quality and accurate line number reporting.",
			len(foundModules),
			utils.Top5StringList(foundModules))
	}

	r.ReportModuleIssue(message, report.IssueSeverityMedium, foundModules)
	r.MakeRecommendation("Ensure you compile the application with debug symbols (PDBs) and include them in the upload.")
}
