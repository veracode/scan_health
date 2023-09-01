package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
)

func thirdParty(r *report.Report) {
	var thirdPartyFilePatterns = []string{
		"7z.dll",
		"7-zip.dll",
		"Google.*.dll",
		"7za.exe",
		"AutoMapper.dll",
		"AutoMapper.*.dll",
		"Azure.*.dll",
		"BouncyCastle.*",
		"Castle.Core.*",
		"Castle.Windsor.*",
		"componentspace.saml2.dll",
		"Dapper.dll",
		"Dapper.*.dll",
		"devexpress.*",
		"entityframework.*",
		"Google.Protobuf.dll",
		"gradle-wrapper.jar",
		"GraphQL.*.dll",
		"itextsharp.dll",
		"log4net.dll",
		"microsoft.*.dll",
		"microsoft.*.pdb",
		"newrelic.*.dll",
		"newtonsoft.json.*",
		"ninject.*.dll",
		"org.eclipse.*.jar",
		"Serilog.dll",
		"syncfusion.*",
		"system.*.dll",
		"Telerik.*.dll",
		"WebGrease.dll",
		"phantomjs.exe",
		"Moq.dll",
		"ComponentSpace.SAML2.dll",
		"*aspnet-codegenerator*",
		"sni.dll",
		"AntiXssLibrary.dll",
		"Antlr3.Runtime.dll",
		"FluentValidation.dll",
	}

	var selectedThirdPartyModules []string

	for index, uploadedFile := range r.UploadedFiles {
		if utils.IsFileNameInFancyList(uploadedFile.Name, thirdPartyFilePatterns) {
			r.UploadedFiles[index].IsThirdParty = true
		}
	}

	for index, module := range r.Modules {
		if utils.IsFileNameInFancyList(module.Name, thirdPartyFilePatterns) {
			r.Modules[index].IsThirdParty = true

			if r.Modules[index].IsSelected() {
				if !utils.IsStringInStringArray(module.Name, selectedThirdPartyModules) {
					selectedThirdPartyModules = append(selectedThirdPartyModules, module.Name)
				}
			}
		}
	}

	if len(selectedThirdPartyModules) == 0 {
		return
	}

	var message = fmt.Sprintf("A third-party component was selected as an entry point: \"%s\".", selectedThirdPartyModules[0])

	if len(selectedThirdPartyModules) > 1 {
		message = fmt.Sprintf("%d third-party components were selected as an entry point: %s.", len(selectedThirdPartyModules), utils.Top5StringList(selectedThirdPartyModules))
	}

	r.ReportModuleIssue(message, report.IssueSeverityMedium, selectedThirdPartyModules)
	r.MakeRecommendation("Only select first party components as the entry points for the analysis. This would typically be any standalone binary or the modules that contain views/controllers.")
}
