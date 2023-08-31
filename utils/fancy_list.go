package utils

import "strings"

// IsFileNameInFancyList tests a wildcard string
// Use a single * for case-insensitive wildcard search, e.g. "*.exe"
// Use ! for case-sensitive wildcard search, e.g. "!*.EXE"
// Use ^ for contains e.g. "^test"
// Otherwise will check for equality (case-insensitive)
func IsFileNameInFancyList(fileName string, fancyList []string) bool {
	preFormattedFileName := strings.TrimSpace(fileName)

	for _, moduleFromFancyList := range fancyList {
		formattedFileName := preFormattedFileName

		// Are we doing a case-insensitive search
		if strings.Count(moduleFromFancyList, "!") == 1 {
			moduleFromFancyList = strings.ReplaceAll(moduleFromFancyList, "!", "")
		} else {
			formattedFileName = strings.ToLower(formattedFileName)
			moduleFromFancyList = strings.ToLower(moduleFromFancyList)
		}

		// Are we doing a contains/within search
		if strings.Count(moduleFromFancyList, "^") == 1 {
			moduleFromFancyList = strings.ReplaceAll(moduleFromFancyList, "^", "")
			if strings.Contains(formattedFileName, moduleFromFancyList) {
				return true
			}

			continue
		}

		// There can only be one * wildcard present
		if strings.Count(moduleFromFancyList, "*") == 1 {
			// At the start (*.xyz)
			if strings.HasPrefix(moduleFromFancyList, "*") {
				if strings.HasSuffix(formattedFileName, strings.ReplaceAll(moduleFromFancyList, "*", "")) {
					return true
				}

				// At the end (abc.*)
			} else if strings.HasSuffix(moduleFromFancyList, "*") {
				if strings.HasPrefix(formattedFileName, strings.ReplaceAll(moduleFromFancyList, "*", "")) {
					return true
				}

				// Or somewhere in the middle (abc.*.xyz)
			} else {
				parts := strings.Split(moduleFromFancyList, "*")

				if strings.HasPrefix(formattedFileName, parts[0]) && strings.HasSuffix(formattedFileName, parts[1]) {
					return true
				}
			}
		} else if strings.EqualFold(formattedFileName, moduleFromFancyList) {
			return true
		}
	}

	return false
}
