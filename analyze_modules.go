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

	if len(modules) > 1000 {
		report.WriteString(fmt.Sprintf(
			"⚠️  %d modules were present. This is a lot of modules which is usually an indicator that something is not correct.\n",
			len(modules)))
	}

	if report.Len() > 0 {
		printTitle("Modules")
		colorPrintf(report.String() + "\n")
	}
}
