package checks

import "github.com/antfie/scan_health/v2/report"

func previousScan(r *report.Report) {
	return

	// TODO

	r.ReportIssue("The uploaded modules for this scan do not match the modules you uploaded for the previous scan. In this scan TODO modules were identified, and TODO were selected for scanning, whereas in the previous scan we observed TODO modules, TODO of which had been selected for scanning. Also noticeable was the total analysis size difference between the two scans.", report.IssueSeverityMedium)
	r.MakeRecommendation("The use of automation will lead to consistent scans.")
}
