package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
	"strings"
)

func repositories(r *report.Report) {
	var foundFiles []string

	for _, uploadedFile := range r.UploadedFiles {
		if strings.EqualFold(uploadedFile.Name, "fsmonitor-watchman.sample") || strings.EqualFold(uploadedFile.Name, "FETCH_HEAD") {
			if !utils.IsStringInStringArray(uploadedFile.Name, foundFiles) {
				foundFiles = append(foundFiles, uploadedFile.Name)
			}
		}
	}

	if len(foundFiles) == 0 {
		return
	}

	issueText := "A git repository had been uploaded."

	if len(foundFiles) > 1 {
		issueText = "Git repositories were uploaded."
	}

	r.ReportFileIssue(fmt.Sprintf("%s Git repositories can contain sensitive information or secrets. We do not recommend uploading repositories as this can increase the time it takes to upload and scan. It can also result in undesirable modules being presented for analysis", issueText), report.IssueSeverityMedium, foundFiles)
	r.MakeRecommendation("Do not upload source code repositories.")
}
