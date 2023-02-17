package main

import (
	"fmt"
	"strings"
)

func (data Data) analyzeModules() {
	var report strings.Builder

	if len(data.PrescanModuleList.Modules) > 1000 {
		report.WriteString(fmt.Sprintf(
			"⚠️  %d modules were identified. This is a lot of modules which is usually an indicator that something is not correct\n",
			len(data.PrescanModuleList.Modules)))
	}

	if len(data.DetailedReport.StaticAnalysis.Modules) > 100 {
		report.WriteString(fmt.Sprintf(
			"⚠️  %d modules were selected for analysis. This is a lot of modules which is usually an indicator that something is not correct\n",
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
		report.WriteString(fmt.Sprintf(
			"⚠️  %d Third party module%s selected: %s\n",
			len(thirdPartyModules),
			pluralise(len(thirdPartyModules)),
			top5StringList(thirdPartyModules)))
	}

	if len(junkModulesSelected) > 0 {
		report.WriteString(fmt.Sprintf(
			"⚠️  %d module%s selected that likely should not be: %s\n",
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

	var errors []string

	for _, module := range data.PrescanModuleList.Modules {
		if module.HasFatalErrors {
			formattedError := fmt.Sprintf("\"%s\": %s", module.Name, module.Status)

			if !isStringInStringArray(formattedError, errors) {
				errors = append(errors, formattedError)
			}
		}
	}

	for _, errors := range errors {
		report.WriteString(fmt.Sprintf("⚠️  %s\n", errors))
	}

	if report.Len() > 0 {
		printTitle("Module Errors")
		colorPrintf(report.String() + "\n")
	}
}

func (data Data) analyzeModuleWarnings() {
	var report strings.Builder

	var warnings []string

	for _, module := range data.PrescanModuleList.Modules {
		if module.IsThirdParty {
			continue
		}

		formattedModuleName := strings.ToLower(module.Name)

		for _, issue := range module.Issues {
			if issue.Details == "No supporting files or PDB files" {
				if strings.HasSuffix(formattedModuleName, ".jar") || strings.HasSuffix(formattedModuleName, ".war") || strings.HasSuffix(formattedModuleName, ".ear") {
					continue
				}

				if strings.HasSuffix(formattedModuleName, ".map") || strings.Contains(formattedModuleName, "_nodemodule_") {
					continue
				}

				if strings.HasPrefix(formattedModuleName, "JS files within") {
					continue
				}
			} else if strings.HasSuffix(formattedModuleName, ".dll") || strings.HasSuffix(formattedModuleName, ".exe") {
				data.makeRecommendation("Ensure you include PDB files for all 1st and 2nd party .NET components. This enables Veracode to accurately report line numbers for any found flaws")
			}

			formattedIssue := fmt.Sprintf("\"%s\": %s", module.Name, issue.Details)

			if !isStringInStringArray(formattedIssue, warnings) {
				warnings = append(warnings, formattedIssue)
			}
		}

		for _, statusMessage := range strings.Split(module.Status, ",") {
			formattedStatusMessage := strings.TrimSpace(statusMessage)

			if strings.HasPrefix(formattedStatusMessage, "Unsupported Framework") {
				// These are captured under the issue details
				continue
			}

			formattedIssue := fmt.Sprintf("\"%s\": %s", module.Name, formattedStatusMessage)

			if !isStringInStringArray(formattedIssue, warnings) {
				warnings = append(warnings, formattedIssue)

				if strings.Contains(formattedStatusMessage, "Missing Supporting Files") {
					data.makeRecommendation("Be sure to include all the components that make up the application within the upload. Do not omit any 2nd or third party libraries from the upload")
				}
			}
		}

	}

	for _, warning := range warnings {
		report.WriteString(fmt.Sprintf("⚠️  %s\n", warning))
	}

	if report.Len() > 0 {
		printTitle("Module Warnings")
		colorPrintf(report.String() + "\n")
	}
}
