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

	detectUnwantedFiles(&report, files, ".zip", "nested zip", "Do not upload archives within the upload package (nested archives)")
	detectUnwantedFiles(&report, files, ".7z", "7-zip", "Veracode does not support 7-zip. Consider zip files instead")
	detectUnwantedFiles(&report, files, ".java", "Java source code", "Do not upload Java source code")
	detectUnwantedFiles(&report, files, ".cs", "C# source code", "Do not upload C# source code")

	if report.Len() > 0 {
		printTitle("Files Uploaded")
		colorPrintf(report.String() + "\n")
	}
}

func detectUnwantedFiles(report *strings.Builder, files []string, suffix, name, advice string) {
	var foundFiles []string

	for _, fileName := range files {
		if strings.HasSuffix(strings.ToLower(fileName), suffix) && !isStringInStringArray(fileName, foundFiles) {
			foundFiles = append(foundFiles, fileName)
		}
	}

	if len(foundFiles) == 0 {
		return
	}

	report.WriteString(fmt.Sprintf(
		"❌ %d %s file%s: %s\n✅ %s\n",
		len(foundFiles),
		name,
		pluralise(len(foundFiles)),
		top5StringList(foundFiles), advice))
}
