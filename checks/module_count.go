package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
)

func moduleCount(r *report.Report) {
	selectedModules := r.GetSelectedModules()

	if len(selectedModules) > utils.MaximumModuleSelectedCountThreshold {
		var selectedModuleNames []string

		for _, module := range selectedModules {
			selectedModuleNames = append(selectedModuleNames, module.Name)
		}

		r.ReportModuleIssue(fmt.Sprintf("%d modules were selected as entry points for analysis. This is a lot of modules and is an indicator that the module selection configuration may be suboptimal.", len(r.GetSelectedModules())), report.IssueSeverityMedium, selectedModuleNames)
		r.MakeRecommendation("Ensure the correct modules have been selected for analysis and that the packaging guidance has been followed.")
		r.MakeRecommendation("Only select the main entry points of the application and not libraries, as documented here: https://community.veracode.com/s/article/What-are-Modules-and-how-do-my-results-change-based-on-what-I-select.")
		r.MakeRecommendation("Consider splitting the application profile into smaller application profiles, one for each deployable component of the application in production. This will make it easier for the various owners of the individual components to filter out areas they are not responsible for, improve scan performance and enable the security team to see the risk of each specific component that makes the whole product family.")
	}

	if len(r.Modules) <= utils.MaximumModuleCountThreshold {
		return
	}

	var modules []string

	for _, module := range r.Modules {
		modules = append(modules, module.Name)
	}

	r.ReportModuleIssue(fmt.Sprintf("%d modules were identified from the upload. This is a lot of modules and is an indicator that the compilation/upload may be suboptimal.", len(r.Modules)), report.IssueSeverityMedium, modules)
	r.MakeRecommendation("Ensure the correct modules have been selected for analysis and that the packaging guidance has been followed.")
}
