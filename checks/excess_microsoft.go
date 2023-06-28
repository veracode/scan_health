package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
)

// Test cases
// https://analysiscenter.veracode.com/auth/index.jsp#StaticOverview::380748:24113946:24085146:24100796::::

func excessMicrosoft(r *report.Report) {
	var filePatterns = []string{
		"csc.exe",
	}

	var foundFiles = r.FancyListMatchUploadedFiles(filePatterns)

	if len(foundFiles) == 0 {
		return
	}

	message := fmt.Sprintf(
		"An unnecessary .NET Roslyn/Runtime component was observed to have been uploaded: \"%s\". Veracode is aware of the .NET runtime, so it is not necessary to include these components. Following the .net packaging will result in the fewest Microsoft libraries in the upload package.",
		foundFiles[0])

	if len(foundFiles) > 1 {
		message = fmt.Sprintf(
			"%d unnecessary .NET Roslyn/Runtime components were observed to have been uploaded: %s. Veracode is aware of the .NET runtime, so it is not necessary to include these components. Following the .net packaging will result in the fewest Microsoft libraries in the upload package.",
			len(foundFiles),
			utils.Top5StringList(foundFiles))
	}

	r.ReportIssue(message, report.IssueSeverityMedium)
	r.MakeRecommendation("Where possible do not include unnecessary Microsoft runtime components.")
}
