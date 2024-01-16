package checks

import "github.com/veracode/scan_health/v2/report"

// Test cases
// https://analysiscenter.veracode.com/auth/index.jsp#StaticOverview::1656378:24066707:24037910:24053560::::

func looseClassFiles(r *report.Report) {
	var filePatterns = []string{
		"*.class",
	}

	var modulePatterns = []string{
		"class files within*",
	}

	var foundFiles = r.FancyListMatchUploadedFiles(filePatterns)
	var foundModules = r.FancyListMatchModules(modulePatterns)

	if len(foundFiles) > 0 {
		r.ReportFileIssue("Java class files were not packaged within a JAR, WAR or EAR file. This is an indicator that the compilation/upload may be suboptimal.", report.IssueSeverityMedium, foundFiles)
		r.MakeRecommendation("Do not upload Java class files. Veracode requires the Java application to be compiled into a JAR, WAR or EAR file as per the packaging instructions.")
	}

	if len(foundModules) > 0 {
		r.ReportModuleIssue("Java class files were not packaged within a JAR, WAR or EAR file. This is an indicator that the compilation/upload may be suboptimal.", report.IssueSeverityMedium, foundModules)
		r.MakeRecommendation("Do not upload Java class files. Veracode requires the Java application to be compiled into a JAR, WAR or EAR file as per the packaging instructions.")
	}
}
