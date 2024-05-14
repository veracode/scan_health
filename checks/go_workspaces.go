package checks

import "github.com/veracode/scan_health/v2/report"

// Test cases
// https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:49327:2134405:36108598:36077125:36092775::::6309710

func unsupportedGoWorkspaceFiles(r *report.Report) {
	var filePatterns = []string{
		"go.work",
		"go.work.sum",
	}

	var foundFiles = r.FancyListMatchUploadedFiles(filePatterns)

	if len(foundFiles) > 0 {
		r.ReportFileIssue("Go workspaces were identified. This is an indicator that the compilation/upload may be suboptimal as Veracode SAST does not currently support go multi-module workspaces.", report.IssueSeverityMedium, foundFiles)
		r.MakeRecommendation("Veracode SAST does not currently support go multi-module workspaces. Please follow the Go packaging instructions: https://docs.veracode.com/r/compilation_go.")
		r.MakeRecommendation("The Veracode CLI can be used to package go apps: https://docs.veracode.com/r/About_auto_packaging")
	}
}
