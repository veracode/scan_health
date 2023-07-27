package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
)

// Test cases
// ==========
// https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:71832:671347:27229752:27200622:27216272::::1670856

func unselectedFirstParty(r *report.Report) {
	var foundModules []string

	for _, module := range r.Modules {
		if !module.IsDependency() && !module.IsIgnored && !module.IsSelected() && !module.IsThirdParty {
			if !utils.IsStringInStringArray(module.Name, foundModules) {
				foundModules = append(foundModules, module.Name)
			}
		}
	}

	if len(foundModules) == 0 {
		return
	}

	message := fmt.Sprintf("A potential first-party module was not selected for analysis: \"%s\".", foundModules[0])

	if len(foundModules) > 1 {
		message = fmt.Sprintf("%d potential first-party modules were not selected for analysis: %s.", len(foundModules), utils.Top5StringList(foundModules))
	}

	r.ReportModuleIssue(message, report.IssueSeverityMedium, foundModules)
	r.MakeRecommendation("Under-selection of first party modules affects results quality. Ensure the correct entry points have been selected as recommended and refer to this article: https://community.veracode.com/s/article/What-are-Modules-and-how-do-my-results-change-based-on-what-I-select.")
}
