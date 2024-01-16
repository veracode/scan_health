package checks

import (
	"fmt"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
)

func tooManyFilesUploaded(r *report.Report) {
	if len(r.UploadedFiles) <= utils.MaximumUploadedFileCountThreshold {
		return
	}

	r.ReportIssue(fmt.Sprintf("Potentially too many files may have been included in the upload (%d). This can result in many modules and long scan times.", len(r.UploadedFiles)), report.IssueSeverityMedium)
}
