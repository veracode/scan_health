package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
)

func missingSupportingFiles(r *report.Report) {
	return

	// TODO

	issueDescription := "Veracode can only scan what has been uploaded. Missing files means reduced scan coverage."

	r.ReportIssue(fmt.Sprintf("There was a module with missing supporting files. %s", issueDescription), report.IssueSeverityMedium)
	r.MakeRecommendation("For optimal scan quality review and resolve the missing supporting files identified in the Review Modules page. To the left of the module name there is an expander button [+] that when pressed will itemize the missing files.")
}
