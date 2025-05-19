package checks

import (
	"fmt"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
)

// Test cases
// https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:78556:1097307:27080380:27051279:27066929::::5141653

func unscannableJava(r *report.Report) {
	var foundModules []string

	for _, module := range r.Modules {
		if module.IsJavaModule() && module.HasFatalErrors() {
			if !utils.IsStringInStringArray(module.Name, foundModules) {
				foundModules = append(foundModules, module.Name)
			}
		}
	}

	if len(foundModules) == 0 {
		return
	}

	message := fmt.Sprintf("A Java module \"%s\" was not scannable as it contained no Java binaries.", foundModules[0])

	if len(foundModules) > 1 {
		message = fmt.Sprintf("%d Java modules were found that contained no Java binaries: %s.", len(foundModules), utils.Top5StringList(foundModules))
	}

	r.ReportModuleIssue(message, report.IssueSeverityHigh, foundModules)
	r.MakeRecommendation("Veracode requires the Java application to be compiled into a JAR, WAR or EAR file as per the packaging instructions.")
	r.MakeRecommendation("The Veracode CLI can be used to package Java apps: https://docs.veracode.com/r/About_auto_packaging.")
}
