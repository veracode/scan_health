package checks

import (
	"fmt"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
)

func flawCount(r *report.Report) {
	if r.Flaws.Total == 0 {
		r.ReportIssue("No flaws were found in this scan. This is usually due to scan misconfiguration.", report.IssueSeverityMedium)
		r.MakeRecommendation("When no flaws have been found this can be an indication that incorrect modules were selected, or the main application was not selected for analysis.")
		return
	}

	if r.Flaws.Total > utils.MaximumFlawCountThreshold {
		r.ReportIssue("A large number of flaws were reported in this scan.", report.IssueSeverityMedium)
		r.MakeRecommendation(fmt.Sprintf("More than %d flaws were found which can be an indication that the scan could be misconfigured.", utils.MaximumFlawCountThreshold))
	}
}
