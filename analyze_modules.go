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
			"⚠️  %d modules were identified. This is a lot of modules which is usually an indicator that something is not correct.\n",
			len(data.PrescanModuleList.Modules)))
	}

	if len(data.DetailedReport.StaticAnalysis.Modules) > 100 {
		report.WriteString(fmt.Sprintf(
			"⚠️  %d modules were selected for analysis. This is a lot of modules which is usually an indicator that something is not correct.\n",
			len(data.DetailedReport.StaticAnalysis.Modules)))
	}

	if report.Len() > 0 {
		printTitle("Modules")
		colorPrintf(report.String() + "\n")
	}
}
