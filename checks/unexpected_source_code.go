package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
)

// Test cases
// https://analysiscenter.veracode.com/auth/index.jsp#StaticOverview::1656378:24066707:24037910:24053560::::

func unexpectedSourceCode(r *report.Report) {
	processSourceCodeFiles(r, []string{"*.java"}, "Java source code file", []string{"Do not upload Java source code files.", "Veracode requires the Java application to be compiled into a JAR, WAR or EAR file as per the packaging instructions: https://docs.veracode.com/r/compilation_java. Consider re-packaging these files such that they contain the actual application and not the source code."})
	processSourceCodeFiles(r, []string{"*.cs"}, "C# source code file", []string{"Do not upload C# source code.", "Veracode requires .NET applications to be compiled with debug symbols into EXE, DLL or NUPKG files as per the packaging instructions: https://docs.veracode.com/r/compilation_net. Consider re-packaging these files such that they contain the actual application and not the source code."})
	processSourceCodeFiles(r, []string{"*.sln"}, ".NET solution file", []string{"Do not upload .NET solution files.", "Veracode requires .NET applications to be compiled with debug symbols into EXE, DLL or NUPKG files as per the packaging instructions: https://docs.veracode.com/r/compilation_net. Consider re-packaging these files such that they contain the actual application and not the source code."})
	processSourceCodeFiles(r, []string{"*.csproj"}, "C# project file", []string{"Do not upload C# project files.", "Veracode requires .NET applications to be compiled with debug symbols into EXE, DLL or NUPKG files as per the packaging instructions: https://docs.veracode.com/r/compilation_net. Consider re-packaging these files such that they contain the actual application and not the source code."})
	processSourceCodeFiles(r, []string{"*.c"}, "C source code file", []string{"Do not upload C source code. Veracode requires the application to be compiled with debug symbols."})
	processSourceCodeFiles(r, []string{"*.cpp"}, "C++ source code file", []string{"Do not upload C++ source code. Veracode requires the application to be compiled with debug symbols."})
	processSourceCodeFiles(r, []string{"*.swift"}, "Swift source code file", []string{"Do not upload Swift source code. Consider compiling these files as per the appropriate iOS packaging guidelines."})
}

func processSourceCodeFiles(r *report.Report, filePatterns []string, fileType string, recommendations []string) {
	var foundFiles = r.FancyListMatchUploadedFiles(filePatterns)

	if len(foundFiles) == 0 {
		return
	}

	message := fmt.Sprintf("A %s was uploaded: \"%s\".", fileType, foundFiles[0])

	if len(foundFiles) > 1 {
		message = fmt.Sprintf(
			"%d %ss were uploaded: %s.",
			len(foundFiles),
			fileType,
			utils.Top5StringList(foundFiles))
	}

	r.ReportFileIssue(fmt.Sprintf("%s Veracode does not scan source code files for this technology. We require the compiled application instead of source code for scanning.", message), report.IssueSeverityHigh, foundFiles)

	for _, recommendation := range recommendations {
		r.MakeRecommendation(recommendation)
	}
}
