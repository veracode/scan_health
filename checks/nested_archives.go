package checks

import "github.com/antfie/scan_health/v2/report"

func nestedArchives(r *report.Report) {
	return

	// TODO

	r.ReportIssue("Some of the uploaded archive files contained nested archives. Veracode does not process nested archives so there may have been some components of this upload that were not analyzed.", report.IssueSeverityHigh)
	r.MakeRecommendation("Ensure you do not upload any nested archives because these will not be scanned.")
}

//detectUnwantedFiles(data, &report, files, []string{"*.zip"}, "nested zip file", []string{"Do not upload archives (nested archives) within the upload package"})
