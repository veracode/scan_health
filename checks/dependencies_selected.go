package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
)

func dependenciesSelected(r *report.Report) {
	var foundModules []string

	for _, module := range r.Modules {
		if module.IsSelected && module.IsDependency {
			foundModules = append(foundModules, module.Name)
		}
	}

	if len(foundModules) == 0 {
		return
	}

	var message = fmt.Sprintf("A dependency \"%s\" was selected as an entry point. This could lead to flaws being raised relating to functionality that may be considered un-reachable or not actionable.", foundModules[0])

	if len(foundModules) > 1 {
		message = fmt.Sprintf("%d dependencies were selected as entry points: %s. This could lead to flaws being raised relating to functionality that may be considered un-reachable or not actionable.", len(foundModules), utils.Top5StringList(foundModules))
	}

	r.ReportIssue(message, report.IssueSeverityMedium)
	r.MakeRecommendation("Only select the main entry points of the application and not libraries, as documented here: https://community.veracode.com/s/article/What-are-Modules-and-how-do-my-results-change-based-on-what-I-select.")
}
