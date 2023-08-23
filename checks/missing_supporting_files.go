package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
	"strconv"
	"strings"
)

// Test cases
// Single missing file: https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:44545:420517:28396113:28366807:28382457::::596193
// Many missing files: https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:44545:389694:28394779:4532802:28381124::::496535

func missingSupportingFiles(r *report.Report) {
	var foundModules []string
	var count = 0

	for _, selectedModule := range r.GetSelectedModules() {
		for _, instance := range selectedModule.Instances {
			for _, issue := range instance.Issues {
				if strings.HasPrefix(issue, "Missing Supporting Files") {
					missingFileCount := strings.Split(issue, " ")[4]

					val, err := strconv.Atoi(missingFileCount)

					if err == nil {
						count += val
					}

					if !utils.IsStringInStringArray(selectedModule.Name, foundModules) {
						foundModules = append(foundModules, selectedModule.Name)
					}
				}
			}
		}
	}

	if count == 0 {
		return
	}

	filePlural := ""

	if count > 1 {
		filePlural = "s"
	}

	var message = fmt.Sprintf("A module \"%s\" was found to be missing %d file%s.", foundModules[0], count, filePlural)

	if len(foundModules) > 1 {
		message = fmt.Sprintf("%d modules were selected as entry points that were found to be missing a total of %d file%s: %s.", len(foundModules), count, filePlural, utils.Top5StringList(foundModules))
	}

	issueDescription := "Veracode can only scan what has been uploaded. Missing files leads to reduced scan coverage."

	r.ReportModuleIssue(fmt.Sprintf("%s %s", message, issueDescription), report.IssueSeverityMedium, foundModules)
	r.MakeRecommendation("For optimal scan quality review and resolve the missing supporting files identified on the Review Modules page. To the left of the module name there is an expander button [+] that when pressed will itemize any missing files.")
}
