package checks

import (
	"fmt"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"strings"
)

// Test cases
// https://analysiscenter.veracode.com/auth/index.jsp#StaticOverview::380748:24113946:24085146:24100796::::

func ignoreJunkFiles(r *report.Report) {
	var filePatternsToIgnore = []string{
		"!LICENSE*",
		".*",
		"*.asmx",
		"*.config",
		"*.cs",
		"*.eot",
		"*.gif",
		"*.ico",
		"*.jpeg",
		"*.jpg",
		"*.less",
		"*.manifest",
		"*.map",
		"*.markdown",
		"*.md",
		"*.pdf",
		"*.png",
		"*.properties",
		"*.scss",
		"*.sh",
		"*.svg",
		"*.ttf",
		"*.txt",
		"*.woff",
		"*.xml",
		"AUTHORS",
		"CHANGELOG",
		"CONTRIBUTORS",
		"Dockerfile",
		"LICENSE",
		"Makefile",
		"README",
		"Thumbs.db",
	}

	var ignoredFiles []string

	for index, uploadedFile := range r.UploadedFiles {
		// Suppress reporting on .PDB files
		if strings.HasSuffix(strings.ToLower(uploadedFile.Name), ".pdb") {
			r.UploadedFiles[index].IsIgnored = true
			continue
		}

		// Suppress reporting on .gitignore and HEAD repository files
		if strings.EqualFold(uploadedFile.Name, ".gitignore") || strings.EqualFold(uploadedFile.Name, "HEAD") {
			r.UploadedFiles[index].IsIgnored = true
			continue
		}

		if utils.IsFileNameInFancyList(uploadedFile.Name, filePatternsToIgnore) {
			r.UploadedFiles[index].IsIgnored = true
			ignoredFiles = append(ignoredFiles, uploadedFile.Name)
		}
	}

	for index, module := range r.Modules {
		if utils.IsFileNameInFancyList(module.Name, filePatternsToIgnore) {
			r.Modules[index].IsIgnored = true
		}
	}

	if len(ignoredFiles) == 0 {
		return
	}

	message := fmt.Sprintf("An unnecessary file was uploaded: \"%s\".", ignoredFiles[0])

	if len(ignoredFiles) > 1 {
		message = fmt.Sprintf("%d unnecessary files were uploaded: %s.", len(ignoredFiles), utils.Top5StringList(ignoredFiles))
	}

	r.ReportFileIssue(message, report.IssueSeverityMedium, ignoredFiles)
	r.MakeRecommendation("Follow the packaging instructions or use the Veracode auto-packager (https://docs.veracode.com/r/About_auto_packaging) to keep the upload as small as possible to improve upload and scan times.")
}
