package main

import (
	"fmt"
	"sort"
	"strings"
)

func (data Data) analyzeUploadedFiles() {
	var report strings.Builder

	var files []string

	for _, file := range data.PrescanFileList.Files {
		files = append(files, file.Name)

	}

	sort.Strings(files[:])

	if len(files) > 10000 {
		report.WriteString(fmt.Sprintf(
			"⚠️  %d files were present. This is a lot of files which is usually an indicator that something is not correct.\n",
			len(files)))
	}

	detectUnwantedFiles(data, &report, files, ".zip", "nested zip", []string{"Do not upload archives (nested archives) within the upload package"})
	detectUnwantedFiles(data, &report, files, ".7z", "7-zip", []string{"Veracode does not support 7-zip. Consider zip files instead"})
	detectUnwantedFiles(data, &report, files, ".java", "Java source code", []string{"Do not upload Java source code", "Veracode requires Java application to be compiled into a .jar, .war or .ear file"})
	detectUnwantedFiles(data, &report, files, ".class", "Java class", []string{"Do not upload Java class files", "Package Java applications into .jar, .war, .ear files"})
	detectUnwantedFiles(data, &report, files, ".cs", "C# source code", []string{"Do not upload C# source code", "Veracode requires the .NET application to be compiled"})
	detectUnwantedFiles(data, &report, files, ".c", "C source code", []string{"Do not upload C source code", "Veracode requires the application to be compiled with debug symbols"})
	detectUnwantedFiles(data, &report, files, ".test.dll", "Test artefacts", []string{"Do not upload any test code"})

	if report.Len() > 0 {
		printTitle("Files Uploaded")
		colorPrintf(report.String() + "\n")
	}
}

func detectUnwantedFiles(data Data, report *strings.Builder, files []string, suffix, name string, recommendations []string) {
	var foundFiles []string

	for _, fileName := range files {
		if strings.HasSuffix(strings.ToLower(fileName), suffix) && !isStringInStringArray(fileName, foundFiles) {
			foundFiles = append(foundFiles, fileName)
		}
	}

	if len(foundFiles) == 0 {
		return
	}

	for _, recommendation := range recommendations {
		data.makeRecommendation(recommendation)
	}

	report.WriteString(fmt.Sprintf(
		"❌ %d %s file%s: %s\n",
		len(foundFiles),
		name,
		pluralise(len(foundFiles)),
		top5StringList(foundFiles)))
}
