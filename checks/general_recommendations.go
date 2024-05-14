package checks

import (
	"fmt"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
)

func generalRecommendations(r *report.Report) {
	if len(r.Recommendations) < 1 {
		return
	}

	if utils.IsStringInStringArray("module", r.Recommendations) {
		r.MakeRecommendation(fmt.Sprintf("Review the module configuration: %s.", r.Scan.ReviewModulesUrl))
		r.MakeRecommendation("Veracode auto-packaging automates the process of packaging for SAST and SCA: https://docs.veracode.com/r/About_auto_packaging")
		r.MakeRecommendation("Follow the packaging guidance for each supported technology present within the application, as documented here: https://docs.veracode.com/r/compilation_packaging. Note there is also a useful cheat sheet which provides bespoke recommendations based off some questions about the application: https://docs.veracode.com/cheatsheet/.")
		r.MakeRecommendation("Read this guidance on modules: https://community.veracode.com/s/article/What-are-Modules-and-how-do-my-results-change-based-on-what-I-select.")
		r.MakeRecommendation("Read about application profile and sandbox best practices: https://community.veracode.com/s/article/application-profile-and-sandbox-best-practices")
		r.MakeRecommendation("Consider scheduling a consultation to review the packaging and module configuration: https://docs.veracode.com/r/t_schedule_consultation.")
	}
}
