package checks

import (
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
)

// Test cases
// https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:75603:793744:20848159:20819949:20835599::::4857698

func missingSCAComponents(r *report.Report) {
	// If SCA data is not available do not perform this check
	if !r.Scan.IsSCADataAvailable {
		return
	}

	// Only if we did not find any SCA components
	if len(r.SCAComponents) > 0 {
		return
	}

	// And only for certain supported technologies
	if !shouldScanIncludeSCAResults(r) {
		return
	}

	r.ReportIssue("There were no SCA results for this scan. This is usually due to scan misconfiguration. Consider including the relevant SCA artefacts as per the packaging instructions.", report.IssueSeverityMedium)
	r.MakeRecommendation("Follow the packaging guidance for each supported technology present within the application, as documented here: https://docs.veracode.com/r/compilation_packaging. Note there is also a useful cheat sheet which provides bespoke recommendations based off some questions about the application: https://docs.veracode.com/cheatsheet/.")
}

func shouldScanIncludeSCAResults(r *report.Report) bool {
	// https://docs.veracode.com/r/Understanding_the_Upload_and_Scan_Language_Support_Matrix
	var scaSupportedFilePatterns = []string{
		"*.dll",
		"*.exe",
		"*.jar",
		"*.apk",
		"*.aab",
		"*.war",
		"*.ear",
		"*.js",
		"*.ts",
		"*.php",
		"*.lock",
		"package-lock.json",
		"npm-shrinkwrap.json",
		"go.sum",
		"vendor.json",
		"*.deps.json",
		"*.py",
	}

	for _, uploadedFile := range r.UploadedFiles {
		if utils.IsFileNameInFancyList(uploadedFile.Name, scaSupportedFilePatterns) {
			return true
		}
	}

	return false
}
