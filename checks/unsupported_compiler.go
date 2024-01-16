package checks

import (
	"fmt"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
)

func unsupportedPlatformOrCompiler(r *report.Report) {
	var foundModules []string

	for _, module := range r.Modules {
		if !module.HasFatalErrors() {
			continue
		}

		// Ignore junk
		if module.IsIgnored || module.IsThirdParty {
			continue
		}

		if module.HasStatus("(Fatal)Unsupported Platform") || module.HasStatus("(Fatal)Unsupported Compiler") {
			if !utils.IsStringInStringArray(module.Name, foundModules) {
				foundModules = append(foundModules, module.Name)
			}
		}
	}

	if len(foundModules) == 0 {
		return
	}

	message := fmt.Sprintf("A module could not be scanned because the platform and/or compiler is unsupported: \"%s\".", foundModules[0])

	if len(foundModules) > 1 {
		message = fmt.Sprintf("%d modules could not be scanned because the platforms and/or compilers are unsupported: %s.", len(foundModules), utils.Top5StringList(foundModules))
	}

	r.ReportModuleIssue(message, report.IssueSeverityHigh, foundModules)
	r.MakeRecommendation("Review the packaging documentation to ensure the vendor and version of the compiler is supported.")
}
