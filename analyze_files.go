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

	detectArchives(&report, files)

	if report.Len() > 0 {
		printTitle("Files Uploaded")
		colorPrintf(report.String() + "\n")
	}
}

func detectArchives(report *strings.Builder, files []string) {
	var foundFiles []string

	for _, fileName := range files {
		if strings.HasSuffix(strings.ToLower(fileName), "zip") && !isStringInStringArray(fileName, foundFiles) {
			foundFiles = append(foundFiles, fileName)
		}
	}

	report.WriteString(fmt.Sprintf("Archives found: %s\n", strings.Join(foundFiles, ", ")))
}
