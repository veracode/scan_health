package checks

import (
	"fmt"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
)

// Test cases
// https://analysiscenter.veracode.com/auth/index.jsp#StaticOverview::1656378:24066707:24037910:24053560::::

func detectUnwantedFiles(r *report.Report) {
	processUnwantedFiles(r, []string{"*.7z"}, "7-zip file", []string{"Veracode does not support 7-zip. Consider using zip files instead."})
	processUnwantedFiles(r, []string{"*.coffee"}, "CoffeeScript file", []string{"CoffeeScript source code files will not be scanned because Veracode does not support CoffeeScript.", "Review the JavaScript/TypeScript packaging instructions: https://docs.veracode.com/r/compilation_jscript."})
	processUnwantedFiles(r, []string{"*.sh", "*.ps", "*.ps1", "*.bat"}, "batch/shell script", []string{"Do not upload batch/shell scripts."})
	processUnwantedFiles(r, []string{"setup.exe", "*setup.exe", "*.msi", "installer.exe", "*installer.exe", "*.msix", "*.appx", "*.msixbundle", ".appxbundle"}, "installer", []string{"Do not upload installers or setup programs."})
	processUnwantedFiles(r, []string{"*.pyd"}, "Python-compiled DLL", []string{"Do not upload Python-compiled DLLs (.pyd files)."})
	processUnwantedFiles(r, []string{"*.pyc"}, "compiled Python file", []string{"Veracode requires Python source code for analysis. Do not upload compiled Python files."})
	processUnwantedFiles(r, []string{"*.deploy"}, "ClickOnce \".deploy\" files", []string{"Veracode does not support ClickOnce deployments."})
}

func processUnwantedFiles(r *report.Report, filePatterns []string, fileType string, recommendations []string) {
	var foundFiles = r.FancyListMatchUploadedFiles(filePatterns)

	if len(foundFiles) == 0 {
		return
	}

	message := fmt.Sprintf("A %s was uploaded: \"%s\".", fileType, foundFiles[0])

	if len(foundFiles) > 1 {
		message = fmt.Sprintf(
			"%d %ss were uploaded: %s.",
			len(foundFiles),
			fileType,
			utils.Top5StringList(foundFiles))
	}

	r.ReportFileIssue(fmt.Sprintf("%s Veracode does not currently support SAST scanning of this file.", message), report.IssueSeverityMedium, foundFiles)

	for _, recommendation := range recommendations {
		r.MakeRecommendation(recommendation)
	}

	r.MakeRecommendation("Follow the packaging instructions to keep the upload as small as possible to improve upload and scan times.")
}
