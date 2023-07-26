package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
)

func nestedArchives(r *report.Report) {
	var foundFiles []string

	for _, uploadedFile := range r.UploadedFiles {
		if uploadedFile.Status == "Archive File Within Another Archive" {
			if !utils.IsStringInStringArray(uploadedFile.Name, foundFiles) {
				foundFiles = append(foundFiles, uploadedFile.Name)
			}
		}
	}

	if len(foundFiles) == 0 {
		return
	}

	message := fmt.Sprintf("A nested archive was uploaded: \"%s\".", foundFiles[0])

	if len(foundFiles) > 1 {
		message = fmt.Sprintf(
			"%d nested archives had been uploaded: %s.",
			len(foundFiles),
			utils.Top5StringList(foundFiles))
	}

	r.ReportFileIssue(fmt.Sprintf("%s Veracode does not process nested archives so there may have been some components of this upload that were not analyzed.", message), report.IssueSeverityHigh, foundFiles)
	r.MakeRecommendation("Ensure you do not upload any nested archives because these will not be scanned.")
}
