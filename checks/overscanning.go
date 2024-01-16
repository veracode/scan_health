package checks

import (
	"fmt"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
)

// Test cases
// https://analysiscenter.veracode.com/auth/index.jsp#StaticOverview::380748:24113946:24085146:24100796::::

func overScanning(r *report.Report) {
	var foundModules []string

	// Over-scanning occurs if there is a module that has been selected for scanning
	// which is also listed as a dependency of any other module which has been selected for scanning.
	for _, module := range r.GetSelectedModules() {
		if module.IsSelected() && len(module.DependencyOf) > 0 {
			for _, otherSelectedModule := range r.GetSelectedModules() {
				// So long as not the same module
				// and that other module is the consumer of this module
				if module.Name != otherSelectedModule.Name && utils.IsStringInStringArray(otherSelectedModule.Name, module.DependencyOf) {
					// So long as not already on the list
					if !utils.IsStringInStringArray(module.Name, foundModules) {
						foundModules = append(foundModules, module.Name)
					}
				}
			}
		}
	}

	if len(foundModules) == 0 {
		return
	}

	var message = fmt.Sprintf("A dependency \"%s\" was incorrectly selected as an entry point for analysis. This is because it was already included in the analysis due to being a dependency of other selected modules. This could lead to indetermanistic results due to the presence of duplicate file names, which could contain different implementations.", foundModules[0])

	if len(foundModules) > 1 {
		message = fmt.Sprintf("%d dependencies were incorrectly selected as entry points for analysis: %s. This is because they were already included in the analysis due to being dependencies of other selected modules. This could lead to indetermanistic results due to the presence of duplicate file names, which could contain different implementations.", len(foundModules), utils.Top5StringList(foundModules))
	}

	r.ReportModuleIssue(message, report.IssueSeverityHigh, foundModules)

	r.MakeRecommendation("For optimal scan quality review and resolve the missing supporting files identified on the Review Modules page. ")
	r.MakeRecommendation("Only select the main entry points of the application and not libraries, as documented here: https://community.veracode.com/s/article/What-are-Modules-and-how-do-my-results-change-based-on-what-I-select. To the left of the module name on the Review Modules page there is an expander button [+] that when pressed will itemize module dependencies.")
}
