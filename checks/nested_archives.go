package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
)

func nestedArchives(r *report.Report) {
	var nestedArchives []string

	for _, uploadedFile := range r.UploadedFiles {
		if uploadedFile.Status == "Archive File Within Another Archive" {
			if !utils.IsStringInStringArray(uploadedFile.Name, nestedArchives) {
				nestedArchives = append(nestedArchives, uploadedFile.Name)
			}
		}
	}

	if len(nestedArchives) == 0 {
		return
	}

	message := fmt.Sprintf("A nested archives was uploaded: \"%s\".", nestedArchives[0])

	if len(nestedArchives) > 1 {
		message = fmt.Sprintf(
			"%d appeared to be nested archives: %s.",
			len(nestedArchives),
			utils.Top5StringList(nestedArchives))
	}

	r.ReportIssue(fmt.Sprintf("%s Veracode does not process nested archives so there may have been some components of this upload that were not analyzed.", message), report.IssueSeverityHigh)
	r.MakeRecommendation("Ensure you do not upload any nested archives because these will not be scanned.")
}
