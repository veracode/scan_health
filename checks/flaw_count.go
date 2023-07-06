package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
)

func flawCount(r *report.Report) {
	if r.Flaws.Total == 0 {
		r.ReportIssue("No flaws were found in this scan.", report.IssueSeverityMedium)
		r.MakeRecommendation("When no flaws have been found this can be an indication that modules have not been selected.")
		return
	}

	if r.Flaws.Total > utils.MaximumFlawCountThreshold {
		r.ReportIssue("A large number of flaws were reported in this scan.", report.IssueSeverityMedium)
		r.MakeRecommendation(fmt.Sprintf("More than %d flaws were found which can be an indication that the scan is misconfigured.", utils.MaximumFlawCountThreshold))
	}
}
