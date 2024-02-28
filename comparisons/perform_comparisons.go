package comparisons

import "github.com/veracode/scan_health/v2/report"

func PerformComparisons(a, b *report.Report) {
	// This is a priority ordered list of checks to perform so the customer sees the most important details first
	reportWarnings(a, b)
	reportCommonalities(a, b)
	reportSummaryDifferences("A", a, b)
	reportSummaryDifferences("B", a, b)

	reportSummary(a, b)
}
