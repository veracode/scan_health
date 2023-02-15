package main

import (
	"fmt"
	"sort"
	"strings"
)

func (data Data) analyzeModules() {
	var report strings.Builder

	var modules []string

	for _, module := range data.PrescanModuleList.Modules {
		modules = append(modules, module.Name)
	}

	sort.Strings(modules[:])

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

	if report.Len() > 0 {
		printTitle("Modules")
		colorPrintf(report.String() + "\n")
	}
}

func (data Data) analyzeModuleWarnings() {
	var report strings.Builder

	var warnings []string

	for _, module := range data.PrescanModuleList.Modules {
		for _, issue := range module.Issues {
			if issue.Details == "No supporting files or PDB files" {
				if strings.HasSuffix(module.Name, ".jar") || strings.HasSuffix(module.Name, ".war") || strings.HasSuffix(module.Name, ".ear") {
					continue
				}

				if strings.HasSuffix(module.Name, ".map") || strings.Contains(module.Name, "_nodemodule_") {
					continue
				}
			}

			formattedIssue := fmt.Sprintf("\"%s\": %s", module.Name, issue.Details)

			if !isStringInStringArray(formattedIssue, warnings) {
				warnings = append(warnings, formattedIssue)
			}
		}

	}

	for _, warning := range warnings {
		report.WriteString(fmt.Sprintf("⚠️  %s\n", warning))
	}

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

	if report.Len() > 0 {
		printTitle("Modules")
		colorPrintf(report.String() + "\n")
	}
}
