package main

import (
	"fmt"
	"strings"
)

func (data Data) reportSelectedModules() {
	var report strings.Builder

	// No point in listing every module if there are loads
	if len(data.DetailedReport.StaticAnalysis.Modules) > 25 {
		return
	}

	for _, module := range data.DetailedReport.StaticAnalysis.Modules {
		report.WriteString(fmt.Sprintf(
			"* %s\n",
			module.Name))

	}

	if report.Len() > 0 {
		printTitle("Modules Selected For Analysis")
		colorPrintf(report.String() + "\n")
	}
}

func (data Data) analyzeModules() {
	var report strings.Builder

	// if data.PrescanModuleList.TotalSize > 1e+9 {
	// 	report.WriteString(formatWarningStringFormat(
	// 		"The size of the modules was %s. This is a very large scan and will likely take a long time to run\n",
	// 		humanize.Bytes(uint64(data.PrescanModuleList.TotalSize))))
	// }

	if len(data.PrescanModuleList.Modules) > 1000 {
		report.WriteString(formatWarningStringFormat(
			"%d modules were identified. This is a lot of modules which is usually an indicator that something is not correct\n",
			len(data.PrescanModuleList.Modules)))
	}

	if len(data.DetailedReport.StaticAnalysis.Modules) > 100 {
		report.WriteString(formatWarningStringFormat(
			"%d modules were selected as entry points for analysis. This is a lot of modules which is usually an indicator that something is not correct\n",
			len(data.DetailedReport.StaticAnalysis.Modules)))
	}

	var thirdPartyModules []string
	var junkModulesSelected []string

	for _, module := range data.DetailedReport.StaticAnalysis.Modules {
		if module.IsThirdParty && !isStringInStringArray(module.Name, thirdPartyModules) {
			thirdPartyModules = append(thirdPartyModules, module.Name)
		}

		if module.IsIgnored && !isStringInStringArray(module.Name, junkModulesSelected) {
			junkModulesSelected = append(junkModulesSelected, module.Name)
		}
	}

	if len(thirdPartyModules) > 0 {
		report.WriteString(formatWarningStringFormat(
			"%d 3rd party module%s selected as an entry point which likely should not be: %s\n",
			len(thirdPartyModules),
			pluralise(len(thirdPartyModules)),
			top5StringList(thirdPartyModules)))

		data.makeRecommendation("We recommend only selecting 1st party components as the entry points for the analysis. This would typically be any standalone binary or the modules that contain views/controllers")
	}

	if len(junkModulesSelected) > 0 {
		report.WriteString(formatWarningStringFormat(
			"%d module%s selected that likely should not be: %s\n",
			len(junkModulesSelected),
			pluralise(len(junkModulesSelected)),
			top5StringList(junkModulesSelected)))
	}

	if report.Len() > 0 {
		printTitle("Modules")
		colorPrintf(report.String() + "\n")
	}
}

func (data Data) analyzeModuleFatalErrors() {
	var report strings.Builder

	errors := make(map[string][]string)

	for _, module := range data.PrescanModuleList.Modules {
		if module.HasFatalErrors {
			reason := module.getFatalReason()
			lowerCaseModuleName := strings.ToLower(module.Name)

			if strings.HasPrefix(reason, "No Scannable Binaries") {
				if strings.HasSuffix(lowerCaseModuleName, ".war") || strings.HasSuffix(lowerCaseModuleName, ".ear") || strings.HasSuffix(lowerCaseModuleName, ".jar") {
					data.makeRecommendation(fmt.Sprintf("Java %s files with no compiled Java classes suggests incorrect packaging and will not be scanned for Java flaws", lowerCaseModuleName))
					data.makeRecommendation("Veracode requires the Java application to be compiled into a .jar, .war or .ear file")
					data.makeRecommendation("Do not upload Java source code files. They will not be scanned")
				}
			}

			if strings.HasSuffix(lowerCaseModuleName, ".docx") {
				// We report on word documents elsewhere
				continue
			}

			if _, isReasonInMap := errors[reason]; !isReasonInMap {
				errors[reason] = []string{}
			}

			if !isStringInStringArray(module.Name, errors[reason]) {
				errors[reason] = append(errors[reason], module.Name)
			}
		}
	}

	for errorMessage, affectedModules := range errors {
		if len(errorMessage) > 1 {
			errorMessage = errorMessage + ": "
		}

		if len(affectedModules) > 1 {
			report.WriteString(formatErrorStringFormat(
				"%dx %s%s\n",
				len(affectedModules),
				errorMessage,
				top5StringList(affectedModules)))
		} else {
			report.WriteString(formatErrorStringFormat(
				"%s%s\n",
				errorMessage,
				top5StringList(affectedModules)))
		}
	}

	if report.Len() > 0 {
		printTitle("Module Errors")
		colorPrintf(report.String() + "\n")
	}
}

func (data Data) analyzeModuleWarnings() {
	var report strings.Builder

	warnings := make(map[string][]string)

	for _, module := range data.PrescanModuleList.Modules {
		if module.HasFatalErrors {
			continue
		}

		if module.IsThirdParty {
			continue
		}

		formattedModuleName := strings.ToLower(module.Name)

		// Look for instances where we are not selecting JS files within
		if strings.Contains(formattedModuleName, "js files within") {
			found := false
			for _, selectedModule := range data.DetailedReport.StaticAnalysis.Modules {
				if strings.EqualFold(selectedModule.Name, module.Name) {
					found = true
				}
			}

			if !found {
				key := "JavaScript module was not selected for analysis"
				warnings[key] = append(warnings[key], module.Name)
				data.makeRecommendation("Veracode extracted JavaScript from the uploaded application. Consider selecting the relevant 'JS files within ...' modules for analysis to cover the JavaScript risk.")
			}
		}

		for _, issue := range module.Issues {
			if strings.HasPrefix(issue.Details, "Unsupported framework") {
				continue
			}

			if strings.Contains(issue.Details, "This application is using Typescript") {
				continue
			}

			if strings.HasPrefix(issue.Details, "Support Issue: ") {
				issue.Details = strings.Replace(issue.Details, "Support Issue: ", "", 1)
			}

			if strings.HasPrefix(issue.Details, "No precompiled files were found for this .NET web application") {
				recommendation := "If this is an ASP.NET application, please precompile the project and upload all generated assemblies"
				issue.Details = strings.Replace(issue.Details, fmt.Sprintf(". %s.", recommendation), "", 1)
				data.makeRecommendation(recommendation)
				data.makeRecommendation("When precompiling ASP.NET WebForms and MVC View ensure you specify the -fixednames flag")
			}

			if strings.Contains(issue.Details, "because we think it is minified") ||
				(strings.HasSuffix(strings.ToLower(issue.Details), ".js") && strings.Contains(issue.Details, "dist/")) {
				data.makeRecommendation("Veracode requires that you submit JavaScript as source code in a format readable by developers. Avoid build steps that minify, obfuscate, bundle, or otherwise compress JavaScript sources")
				data.makeRecommendation("Do not upload files that are concatenated or minified. Veracode ignores files that have filenames that suggest that they are concatenated or minified")
				data.makeRecommendation("Review the JavaScript/TypeScript packaging cheatsheet: https://nhinv11.github.io/#/JavaScript%20/%20TypeScript")
				data.makeRecommendation("Consider using the unofficial JavaScript/TypeScript packaging tool: https://github.com/fw10/veracode-javascript-packager")
			}

			if strings.Contains(issue.Details, "test/") {
				data.makeRecommendation("Do not upload any testing artefacts")
			}

			if strings.Contains(module.Name, "Class files within") {
				issue.Details = ""
				data.makeRecommendation("Veracode requires Java application to be compiled into a .jar, .war or .ear file")
			}

			if issue.Details == "No supporting files or PDB files" {
				if strings.HasSuffix(formattedModuleName, ".dll") || strings.HasSuffix(formattedModuleName, ".exe") {
					data.makeRecommendation("Ensure you include PDB files for all 1st and 2nd party .NET components. This enables Veracode to accurately report line numbers for any found flaws")
				} else {
					continue
				}
			}

			if _, isMessageInMap := warnings[issue.Details]; !isMessageInMap {
				warnings[issue.Details] = []string{}
			}

			if !isStringInStringArray(module.Name, warnings[issue.Details]) {
				warnings[issue.Details] = append(warnings[issue.Details], module.Name)
			}
		}

		for _, statusMessage := range strings.Split(module.Status, ",") {
			if module.Status == "OK" {
				continue
			}

			formattedStatusMessage := strings.TrimSpace(statusMessage)

			if strings.HasPrefix(formattedStatusMessage, "Unsupported Framework") {
				// Ignore this
				continue
			}

			if strings.HasPrefix(formattedStatusMessage, "Support Issue") {
				// Ignore this
				continue
			}

			if strings.HasPrefix(formattedStatusMessage, "PDB Files Missing") {
				if strings.HasSuffix(formattedModuleName, ".dll") || strings.HasSuffix(formattedModuleName, ".exe") {
					formattedStatusMessage = "Modules with missing PDB files"
					data.makeRecommendation("Ensure you include PDB files for all 1st and 2nd party .NET components. This enables Veracode to accurately report line numbers for any found flaws")
				} else {
					continue
				}
			}

			if strings.Contains(formattedStatusMessage, "Missing Supporting Files") {
				formattedStatusMessage = "Modules with missing supporting files"
				data.makeRecommendation("Be sure to include all the components that make up the application within the upload. Do not omit any 2nd or third party libraries from the upload")
			}

			if !isStringInStringArray(module.Name, warnings[formattedStatusMessage]) {
				warnings[formattedStatusMessage] = append(warnings[formattedStatusMessage], module.Name)
			}
		}
	}

	for warningMessage, affectedModules := range warnings {
		if len(warningMessage) > 1 {
			warningMessage = warningMessage + ": "
		}

		if len(affectedModules) > 1 {
			report.WriteString(formatWarningStringFormat(
				"%dx %s%s\n",
				len(affectedModules),
				warningMessage,
				top5StringList(affectedModules)))
		} else {
			report.WriteString(formatWarningStringFormat(
				"%s%s\n",
				warningMessage,
				top5StringList(affectedModules)))
		}
	}

	if report.Len() > 0 {
		printTitle("Module Warnings")
		colorPrintf(report.String() + "\n")
	}
}
