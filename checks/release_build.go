package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
	"strings"
)

func releaseBuild(r *report.Report) {
	detectDotNetReleasePathsInModuleIssues(r)
}

func detectDotNetReleasePathsInModuleIssues(r *report.Report) {
	var foundModules []string

	for _, module := range r.Modules {
		// Only applicable for .net modules
		if !module.IsDotNetOrCPPModule() {
			continue
		}

		for _, issue := range module.Issues {
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

	r.ReportIssue(fmt.Sprintf("%s Unit tests and mocks can make it difficult to select the correct application entry points for analysis. This is because for most cases Veracode permits users to select only the components for analysis that are not themselves depended upon by other components within the upload. Furthermore, scanning unit tests will surface flaws that will not be present in a production environment and commonly they contain hard-coded credentials for testing purposes.", message), report.IssueSeverityMedium)
	r.MakeRecommendation("Do not upload any testing artifacts unless they will go into the production environment.")
}
