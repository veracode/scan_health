package checks

import (
	"fmt"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
)

func duplicateModules(r *report.Report) {
	sameDuplicates, differentDuplicates := detectDuplicates(r)

	if len(differentDuplicates) > 0 {
		var duplicates []string

		for fileName, count := range differentDuplicates {
			for i := 1; i <= count; i++ {
				duplicates = append(duplicates, fileName)
			}
		}

		message := fmt.Sprintf("A duplicate file name was uploaded but the file hashes were different: %s. This can affect the quality of the scan, result in scans taking longer than expected and lead to indeterministic flaws being raised. This can also cause confusion when interpreting the results. Furthermore, if the scanner found no risk in the first file, risk could be missed in the second file because the scanner only analyses the first filename it comes across when we encounter duplicate files.", utils.Top5StringList(duplicates))

		if len(differentDuplicates) > 1 {
			message = fmt.Sprintf("%d duplicate file names were uploaded but the file hashes were different: %s. This can affect the quality of the scan, result in scans taking longer than expected and lead to indeterministic flaws being raised. This can also cause confusion when interpreting the results. Furthermore, if the scanner found no risk in the first file, risk could be missed in the second file because the scanner only analyses the first filename it comes across when we encounter duplicate files.", len(differentDuplicates), utils.Top5StringList(duplicates))
		}

		r.ReportFileIssue(message, report.IssueSeverityHigh, duplicates)
	}

	if len(sameDuplicates) == 1 {
		for key, value := range sameDuplicates {
			r.ReportFileIssue(fmt.Sprintf("%d duplicates of the file \"%s\" were uploaded. This can affect result in scans taking longer than expected.", value, key), report.IssueSeverityMedium, nil)
		}

	} else if len(sameDuplicates) > 0 {
		r.ReportFileIssue(fmt.Sprintf("%d duplicate files were uploaded. This can affect result in scans taking longer than expected.", len(sameDuplicates)), report.IssueSeverityMedium, nil)
	}

	if len(differentDuplicates)+len(differentDuplicates) > 0 {
		r.MakeRecommendation("De-duplicate the modules/components before upload for optimal scan quality. Typically following the packaging instructions or using the Veracode auto-packager (https://docs.veracode.com/r/About_auto_packaging) will result in an upload with no/few duplicates. Having duplicate file names with different contents can lead to indeterminate scan results.")
		r.MakeRecommendation("Ensure you only upload one version of an application/component of your application in each scan. Refer to considerations on how to scope and structure application profile uploads properly at https://docs.veracode.com/r/request_profile.")
	}
}

func detectDuplicates(r *report.Report) (map[string]int, map[string]int) {
	var sameDuplicates = make(map[string]int)
	var differentDuplicates = make(map[string]int)
	var processedFiles []string

	for _, file := range r.UploadedFiles {
		// We only want to process each unique file name once
		if utils.IsStringInStringArray(file.Name, processedFiles) {
			continue
		}

		if file.IsIgnored || file.IsThirdParty {
			continue
		}

		var uniqueMD5Hashes []string
		var fileCount = 0

		for _, otherFile := range r.UploadedFiles {
			if file.Name == otherFile.Name {
				fileCount++
				if !utils.IsStringInStringArray(otherFile.MD5, uniqueMD5Hashes) {
					uniqueMD5Hashes = append(uniqueMD5Hashes, otherFile.MD5)
				}
			}
		}

		if fileCount > 1 {
			if len(uniqueMD5Hashes) == 1 {
				sameDuplicates[file.Name] = fileCount
			} else {
				differentDuplicates[file.Name] = fileCount
			}
		}

		processedFiles = append(processedFiles, file.Name)
	}

	return sameDuplicates, differentDuplicates
}
