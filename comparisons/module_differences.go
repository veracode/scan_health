package comparisons

import (
	"fmt"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"strings"
)

func reportModuleDifferences(a, b *report.Report) {
	var r strings.Builder

	compareModuleDifferences(&r, a, b)

	if r.Len() > 0 {
		utils.PrintTitle("Module Differences (Ignoring any duplicates)")
		utils.ColorPrintf(r.String())
	}
}

func compareModuleDifferences(r *strings.Builder, a, b *report.Report) {
	var scanANonDuplicatedFiles = getNonDuplicatedFileNames(a.UploadedFiles)
	var scanBNonDuplicatedFiles = getNonDuplicatedFileNames(b.UploadedFiles)

	for _, thisFile := range a.UploadedFiles {
		// Ignore duplicates from scan A
		if !utils.IsStringInStringArray(thisFile.Name, scanANonDuplicatedFiles) {
			continue
		}

		// Ignore duplicates from scan B
		if !utils.IsStringInStringArray(thisFile.Name, scanBNonDuplicatedFiles) {
			continue
		}

		for _, otherFile := range b.UploadedFiles {
			if thisFile.Name == otherFile.Name {
				if thisFile.MD5 != otherFile.MD5 {
					r.WriteString(
						fmt.Sprintf("\"%s\" %s: MD5 = %s, %s: MD5 = %s \n",
							thisFile.Name,
							utils.GetFormattedSideString("A"),
							thisFile.MD5,
							utils.GetFormattedSideString("B"),
							otherFile.MD5))
				}
			}
		}
	}

}

func getNonDuplicatedFileNames(uploadedFiles []report.UploadedFile) []string {
	var duplicateFiles []string
	var processedFiles []string

	for _, file := range uploadedFiles {
		if utils.IsStringInStringArray(file.Name, processedFiles) && !utils.IsStringInStringArray(file.Name, duplicateFiles) {
			duplicateFiles = append(duplicateFiles, file.Name)
		}

		processedFiles = append(processedFiles, file.Name)
	}

	var nonDuplicatedFiles []string

	for _, file := range uploadedFiles {
		if !utils.IsStringInStringArray(file.Name, duplicateFiles) {
			nonDuplicatedFiles = append(nonDuplicatedFiles, file.Name)
		}
	}

	return nonDuplicatedFiles
}
