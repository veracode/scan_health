package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
	"time"
)

// Test cases
// https://analysiscenter.veracode.com/auth/index.jsp#StaticOverview::1656378:24066707:24037910:24053560::::

func regularScans(r *report.Report) {
	// No LastAppActivity data
	if r.LastAppActivity.Year() < 2000 {
		return
	}

	if time.Since(r.LastAppActivity).Hours()/24 <= utils.NotUsingAutomationIfScanOlderThanDays {
		return
	}

	r.ReportIssue(fmt.Sprintf("There have not been recent scans of this application. The application was last scanned on %s which was %s ago. It is not uncommon for new flaws to be reported over time because Veracode is continuously improving their products, and because new SCA vulnerabilities are reported every day, and this could impact the application.", r.LastAppActivity, utils.FormatHumanDurationDays(time.Since(r.LastAppActivity))), report.IssueSeverityMedium)
	r.MakeRecommendation("Regular scanning, preferably via automation will allow the application team to respond faster to any new issues reported.")
}
