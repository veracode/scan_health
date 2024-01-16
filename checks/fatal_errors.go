package checks

import (
	"fmt"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
)

func fatalErrors(r *report.Report) {
	fatalMissingPrimaryDebugSymbols(r)
	fatalNoScannableJavaBinaries(r)
	fatalNestedJars(r)
}

func fatalMissingPrimaryDebugSymbols(r *report.Report) {
	var foundModules []string

	for _, module := range r.Modules {
		if !module.HasFatalErrors() {
			continue
		}

		// Only applicable for .net modules
		if !module.IsDotNetOrCPPModule() {
			continue
		}

		// Ignore junk
		if module.IsIgnored || module.IsThirdParty {
			continue
		}

		if module.HasStatus("Primary Files Compiled without Debug Symbols") {
			if !utils.IsStringInStringArray(module.Name, foundModules) {
				foundModules = append(foundModules, module.Name)
			}
		}
	}

	if len(foundModules) == 0 {
		return
	}

	message := fmt.Sprintf("A module could not be scanned because the first party component did not have debug symbols (e.g., PDB files): \"%s\".", foundModules[0])

	if len(foundModules) > 1 {
		message = fmt.Sprintf("%d modules could not be scanned because the first party components did not have debug symbols (e.g., PDB files): %s.", len(foundModules), utils.Top5StringList(foundModules))
	}

	r.ReportModuleIssue(message, report.IssueSeverityHigh, foundModules)
	r.MakeRecommendation("Include PDB files for as many components as possible, especially first and second party components. This enables Veracode to accurately report line numbers for any flaws found within these components.")
}

func fatalNoScannableJavaBinaries(r *report.Report) {
	var foundModules []string

	for _, module := range r.Modules {
		if !module.HasFatalErrors() {
			continue
		}

		// Only applicable for Java modules
		if !module.IsJavaModule() {
			continue
		}

		if module.HasStatus("No Scannable Binaries") {
			if !utils.IsStringInStringArray(module.Name, foundModules) {
				foundModules = append(foundModules, module.Name)
			}
		}
	}

	if len(foundModules) == 0 {
		return
	}

	message := fmt.Sprintf("A Java module was identified that contained no compiled Java classes: \"%s\". Sometimes this can indicate the file contained only Java source code which is not processed by Veracode.", foundModules[0])

	if len(foundModules) > 1 {
		message = fmt.Sprintf("%d Java modules were identified that contained no compiled Java classes: %s. Sometimes this can indicate the file contained only Java source code which is not processed by Veracode.", len(foundModules), utils.Top5StringList(foundModules))
	}

	r.ReportModuleIssue(message, report.IssueSeverityHigh, foundModules)
	r.MakeRecommendation("Veracode requires the Java application to be compiled into a JAR, WAR or EAR file as per the packaging instructions.")
}

func fatalNestedJars(r *report.Report) {
	var foundModules []string

	for _, module := range r.Modules {
		if !module.HasFatalErrors() {
			continue
		}

		// Only applicable for Java modules
		if !module.IsJavaModule() {
			continue
		}

		if module.HasStatus("does not support jar files nested inside") {
			if !utils.IsStringInStringArray(module.Name, foundModules) {
				foundModules = append(foundModules, module.Name)
			}
		}
	}

	if len(foundModules) == 0 {
		return
	}

	message := fmt.Sprintf("A Java module was identified that contained nested/shaded second or third party components: \"%s\".", foundModules[0])

	if len(foundModules) > 1 {
		message = fmt.Sprintf("%d Java modules were identified that contained nested/shaded second or third party components: %s.", len(foundModules), utils.Top5StringList(foundModules))
	}

	r.ReportModuleIssue(message, report.IssueSeverityHigh, foundModules)
	r.MakeRecommendation("Veracode Static Analysis does not support JAR files nested inside other JAR files, except for Spring Boot applications. Veracode does support analysis of uber-jar files created by the Maven Shade plugin.")
}
