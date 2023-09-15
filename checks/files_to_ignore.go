package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
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
	}

	var ignoredFiles []string

	for index, uploadedFile := range r.UploadedFiles {
		if utils.IsFileNameInFancyList(uploadedFile.Name, filePatternsToIgnore) {
			// Deal with files named like this: .NETCoreApp_Version_v3.1.AssemblyAttributes.cs
			if strings.HasSuffix(strings.ToLower(uploadedFile.Name), ".cs") {
				continue
			}

			// Suppress reporting on .gitignore files
			if strings.EqualFold(uploadedFile.Name, ".gitignore") {
				continue
			}

			r.UploadedFiles[index].IsIgnored = true
			ignoredFiles = append(ignoredFiles, uploadedFile.Name)
		}

		// Ignore .PDB files
		if strings.HasSuffix(strings.ToLower(uploadedFile.Name), ".pdb") {
			r.UploadedFiles[index].IsIgnored = true
		}
	}

	for index, module := range r.Modules {
		if utils.IsFileNameInFancyList(module.Name, filePatternsToIgnore) {
			// Deal with modules named like this: .NETCoreApp_Version_v3.1.AssemblyAttributes.cs
			if strings.HasSuffix(strings.ToLower(module.Name), ".cs") {
				continue
			}

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
	r.MakeRecommendation("Follow the packaging instructions to keep the upload as small as possible to improve upload and scan times.")
}
