package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
	"strings"
)

// Test cases
// https://analysiscenter.veracode.com/auth/index.jsp#StaticOverview::1656378:24066707:24037910:24053560::::

func testingArtefacts(r *report.Report) {
	var testFilePatterns = []string{
		"nunit.framework.dll",
		"Moq.dll",
		"*.test.dll", "*.unittests.dll", "*.unittest.dll", "^mock", "^unittest", "^harness",
	}

	detectSelectedTestingModules(r, testFilePatterns)
	detectTestArtefactsInFileUploads(r, testFilePatterns)
	detectTestArtefactsInModuleIssues(r)
}

func detectSelectedTestingModules(r *report.Report, testFilePatterns []string) {
	var selectedTestingModules = r.FancyListMatchSelectedModules(testFilePatterns)

	if len(selectedTestingModules) == 0 {
		return
	}

	message := fmt.Sprintf("A testing artefact was uploaded and selected as a module: \"%s\".", selectedTestingModules[0])

	if len(selectedTestingModules) > 1 {
		message = fmt.Sprintf(
			"%d testing artefacts were uploaded and selected as a module: %s.",
			len(selectedTestingModules),
			utils.Top5StringList(selectedTestingModules))
	}

	r.ReportModuleIssue(fmt.Sprintf("%s Unit tests and mocks can make it difficult to select the correct application entry points for analysis. This is because for most cases Veracode permits users to select only the components for analysis that are not themselves depended upon by other components within the upload. Furthermore, scanning unit tests will surface flaws that will not be present in a production environment and commonly they contain hard-coded credentials for testing purposes.", message), report.IssueSeverityHigh, selectedTestingModules)
	r.MakeRecommendation("Do not upload any testing artifacts unless they will go into the production environment.")
	r.MakeRecommendation("Do not select any testing artefacts as entry points for analysis.")
}

func detectTestArtefactsInFileUploads(r *report.Report, testFilePatterns []string) {
	var foundFiles = r.FancyListMatchUploadedFiles(testFilePatterns)

	if len(foundFiles) == 0 {
		return
	}

	message := fmt.Sprintf("A testing artefact was uploaded: \"%s\".", foundFiles[0])

	if len(foundFiles) > 1 {
		message = fmt.Sprintf(
			"%d testing artefacts were uploaded: %s.",
			len(foundFiles),
			utils.Top5StringList(foundFiles))
	}

	r.ReportFileIssue(fmt.Sprintf("%s Unit tests and mocks can make it difficult to select the correct application entry points for analysis. This is because for most cases Veracode permits users to select only the components for analysis that are not themselves depended upon by other components within the upload. Furthermore, scanning unit tests will surface flaws that will not be present in a production environment and commonly they contain hard-coded credentials for testing purposes.", message), report.IssueSeverityMedium, foundFiles)
	r.MakeRecommendation("Do not upload any testing artifacts unless they will go into the production environment.")
}

func detectTestArtefactsInModuleIssues(r *report.Report) {
	var foundModules []string

	for _, module := range r.Modules {
		for _, issue := range module.GetAllIssues() {
			if strings.Contains(strings.ToLower(issue), "test/") {
				if !utils.IsStringInStringArray(module.Name, foundModules) {
					foundModules = append(foundModules, module.Name)
				}
			}
		}
	}

	if len(foundModules) == 0 {
		return
	}

	message := fmt.Sprintf("A module appeared to contain or depend on testing artefacts: \"%s\".", foundModules[0])

	if len(foundModules) > 1 {
		message = fmt.Sprintf(
			"%d modules appeared to contain or depend on testing artefacts: %s.",
			len(foundModules),
			utils.Top5StringList(foundModules))
	}

	r.ReportModuleIssue(fmt.Sprintf("%s Unit tests and mocks can make it difficult to select the correct application entry points for analysis. This is because for most cases Veracode permits users to select only the components for analysis that are not themselves depended upon by other components within the upload. Furthermore, scanning unit tests will surface flaws that will not be present in a production environment and commonly they contain hard-coded credentials for testing purposes.", message), report.IssueSeverityMedium, foundModules)
	r.MakeRecommendation("Do not upload any testing artifacts unless they will go into the production environment.")
}
