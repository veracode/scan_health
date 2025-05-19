package checks

import (
	"fmt"
	"strings"

	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
)

// Test cases
// https://analysiscenter.veracode.com/auth/index.jsp#StaticOverview::906008:24693294:24664441:24680091::::

func minifiedJavaScript(r *report.Report) {
	detectMinifiedJSUploads(r)
	detectMinifiedJSModules(r)
}

func detectMinifiedJSUploads(r *report.Report) {
	var foundFiles []string

	for _, uploadedFile := range r.UploadedFiles {
		if strings.HasSuffix(strings.ToLower(uploadedFile.Name), ".min.js") {
			if !utils.IsStringInStringArray(uploadedFile.Name, foundFiles) {
				foundFiles = append(foundFiles, uploadedFile.Name)
			}
		}
	}

	if len(foundFiles) == 0 {
		return
	}

	message := fmt.Sprintf("A minified JavaScript file was uploaded: \"%s\". This file will not be scanned.", foundFiles[0])

	if len(foundFiles) > 1 {
		message = fmt.Sprintf("%d minified JavaScript files were uploaded: %s. These files will not be scanned.", len(foundFiles), utils.Top5StringList(foundFiles))
	}

	r.ReportFileIssue(message, report.IssueSeverityMedium, foundFiles)
	makeJSRecommendations(r)
}

func detectMinifiedJSModules(r *report.Report) {
	var foundModules []string

	for _, module := range r.Modules {
		// Only applicable for JavaScript modules
		if !module.IsJavaScriptModule() {
			continue
		}

		for _, issue := range module.Issues() {
			if strings.Contains(issue, "because we think it is minified") || strings.Contains(strings.ToLower(issue), "dist/") {
				if !utils.IsStringInStringArray(module.Name, foundModules) {
					foundModules = append(foundModules, module.Name)
				}
			}
		}
	}

	if len(foundModules) == 0 {
		return
	}

	message := fmt.Sprintf("A minified JavaScript file was found within this module: \"%s\". This file might not be scanned.", foundModules[0])

	if len(foundModules) > 1 {
		message = fmt.Sprintf("%d minified JavaScript files were found within this module: %s. These files might not be scanned.", len(foundModules), utils.Top5StringList(foundModules))
	}

	r.ReportModuleIssue(message, report.IssueSeverityMedium, foundModules)
	makeJSRecommendations(r)
}

func makeJSRecommendations(r *report.Report) {
	r.MakeRecommendation("Veracode requires that you submit JavaScript as source code in a format readable by developers. Avoid build steps that minify, concatenate, obfuscate, bundle, or otherwise compress JavaScript sources. Veracode ignores files that have filenames that suggest that they are concatenated or minified.")
	r.MakeRecommendation("Review the JavaScript/TypeScript packaging instructions: https://docs.veracode.com/r/compilation_jscript.")
	r.MakeRecommendation("The Veracode CLI can be used to package JavaScript apps: https://docs.veracode.com/r/About_auto_packaging.")
}
