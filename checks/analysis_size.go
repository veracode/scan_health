package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
	"github.com/dustin/go-humanize"
)

func analysisSize(r *report.Report) {
	if r.Scan.AnalysisSize <= utils.MaximumAnalysisSieBytesThreshold {
		return
	}

	r.ReportIssue(fmt.Sprintf("The analysis size of the scan was %s. This is a very large size and will likely take a long time to upload and run.", humanize.Bytes(r.Scan.AnalysisSize)), report.IssueSeverityMedium)
	r.MakeRecommendation("Ensure the correct modules have been selected for analysis and that the packaging guidance has been followed.")
}
