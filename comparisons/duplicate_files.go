package comparisons

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"strings"
)

func reportDuplicateFiles(side string, files []report.UploadedFile) {
	var r strings.Builder
	var processedFiles []string

	for _, thisFile := range files {
		if utils.IsStringInStringArray(thisFile.Name, processedFiles) {
			continue
		}

		md5s := []string{thisFile.MD5}
		var count = 0

		for _, otherFile := range files {
			if thisFile.Name == otherFile.Name {
				count++
				if !utils.IsStringInStringArray(otherFile.MD5, md5s) {
					md5s = append(md5s, otherFile.MD5)
				}
			}
		}

		if len(md5s) > 1 {
			if count == len(md5s) {
				r.WriteString(fmt.Sprintf("\"%s\": %d occurances each with different MD5 hashes\n", thisFile.Name, count))
			} else {
				r.WriteString(fmt.Sprintf("\"%s\": %d occurances with %d different MD5 hashes\n", thisFile.Name, count, len(md5s)))
			}
		}

		processedFiles = append(processedFiles, thisFile.Name)
	}

	if r.Len() > 0 {
		utils.ColorPrintf(utils.GetFormattedSideStringWithMessage(side, fmt.Sprintf("\nDuplicate Files Within Scan %s\n", side)))
		fmt.Print("=============================\n")
		color.HiYellow(r.String())
	}
}
