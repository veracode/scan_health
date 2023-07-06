package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
	"strings"
)

// Test cases
// https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:90561:1490961:27041493:27012395:27028045::::4709647
// https://analysiscenter.veracode.com/auth/index.jsp#StaticOverview::380748:24113946:24085146:24100796::::

func missingPrecompiledFiles(r *report.Report) {
	detectMissingPrecompiledFilesFromUploadedFiles(r)
	detectMissingPrecompiledFilesFromModules(r)
}

func detectMissingPrecompiledFilesFromUploadedFiles(r *report.Report) {
	var templateFileList = []string{
		"*.cshtml",
		"*.ascx",
		"*.aspx",
		"*.asax",
	}

	var foundFiles = r.FancyListMatchUploadedFiles(templateFileList)

	if len(foundFiles) == 0 {
		return
	}

	message := fmt.Sprintf("A .NET view/template/control file was uploaded: %s", foundFiles[0])

	if len(foundFiles) > 1 {
		message = fmt.Sprintf("%d .NET view/template/control files were uploaded: %s.", len(foundFiles), utils.Top5StringList(foundFiles))
	}

	r.ReportIssue(message, report.IssueSeverityHigh)
	recommendPrecompile(r)
}

func detectMissingPrecompiledFilesFromModules(r *report.Report) {
	var foundModules []string

	for _, module := range r.Modules {
		// Only applicable for .net modules
		if !module.IsDotNetOrCPPModule() {
			continue
		}

		// Ignore junk
		if module.IsIgnored || module.IsThirdParty {
			continue
		}

		for _, issue := range module.Issues {
			if strings.Contains(issue, "No precompiled files were found for this .NET web application") {
				if !utils.IsStringInStringArray(module.Name, foundModules) {
					foundModules = append(foundModules, module.Name)
				}
			}
		}
	}

	if len(foundModules) == 0 {
		return
	}

	message := fmt.Sprintf("A .NET component was identified to not contain precompiled files (views/templates/controls): %s.", foundModules[0])

	if len(foundModules) > 1 {
		message = fmt.Sprintf("%d .NET components were identified to not contain precompiled files (view/template/control): %s.", len(foundModules), utils.Top5StringList(foundModules))
	}

	r.ReportIssue(message, report.IssueSeverityMedium)
	recommendPrecompile(r)
}

func recommendPrecompile(r *report.Report) {
	r.MakeRecommendation("If this is an ASP.NET application, please precompile the project and upload all generated assemblies. When precompiling ASP.NET, WebForms and MVC views ensure you specify the \"-fixednames\" flag during compilation. Please pre-compile the project and upload all generated assemblies. If you do not submit precompiled forms, the scan can produce incomplete or incorrect results. Review the documentation on how to do this here: https://docs.veracode.com/r/c_precomp_VS.")
}
