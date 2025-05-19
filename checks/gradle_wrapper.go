package checks

import (
	"github.com/veracode/scan_health/v2/report"
)

// Test cases
// https://analysiscenter.veracode.com/auth/index.jsp#StaticOverview::1656378:24066707:24037910:24053560::::

func gradleWrapper(r *report.Report) {
	var filePatterns = []string{
		"gradle-wrapper.jar",
	}

	var selectedModules = r.FancyListMatchSelectedModules(filePatterns)

	if len(selectedModules) == 0 {
		return
	}

	message := "The only module selected on this scan was \"gradle-wrapper.jar\". This is a known third-party build tool and not the main application to analyse."

	if len(selectedModules) > 1 {
		message = "The \"gradle-wrapper.jar\" build tool selected on this scan for analysis. This is a known third-party component and not the main application to analyse."
	}

	r.ReportModuleIssue(message, report.IssueSeverityHigh, selectedModules)
	r.MakeRecommendation("The \"gradle-wrapper.jar\" component should not be uploaded or selected for analysis.")
	r.MakeRecommendation("The Veracode CLI can be used to package Java apps: https://docs.veracode.com/r/About_auto_packaging.")
}
