package comparisons

import (
	"fmt"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"strings"
)

func reportSummary(a, b *report.Report) {
	var report strings.Builder

	if a.Scan.SubmittedDate.Before(b.Scan.SubmittedDate) {
		report.WriteString(fmt.Sprintf("%s was submitted %s after %s\n", utils.GetFormattedSideString("B"), utils.FormatDuration(b.Scan.SubmittedDate.Sub(a.Scan.SubmittedDate)), utils.GetFormattedSideString("A")))
	} else if a.Scan.SubmittedDate.After(b.Scan.SubmittedDate) {
		report.WriteString(fmt.Sprintf("%s was submitted %s after %s\n", utils.GetFormattedSideString("A"), utils.FormatDuration(a.Scan.SubmittedDate.Sub(b.Scan.SubmittedDate)), utils.GetFormattedSideString("B")))
	}

	aDuration := utils.FormatDuration(a.Scan.Duration)
	bDuration := utils.FormatDuration(b.Scan.Duration)

	if aDuration > bDuration {
		report.WriteString(fmt.Sprintf("%s took longer by %s\n", utils.GetFormattedSideString("A"), utils.FormatDuration(a.Scan.Duration-b.Scan.Duration)))
	} else if aDuration < bDuration {
		report.WriteString(fmt.Sprintf("%s took longer by %s\n", utils.GetFormattedSideString("B"), utils.FormatDuration(b.Scan.Duration-a.Scan.Duration)))
	}

	if report.Len() > 0 {
		utils.PrintTitle("Summary")
		utils.ColorPrintf(report.String())
	}
}
