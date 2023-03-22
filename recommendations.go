package main

import (
	"fmt"
	"strings"
)

func (data Data) outputRecommendations(region string) {
	if (len(*data.Recommendations)) == 0 {
		return
	}

	allRecommendations := strings.Join(*data.Recommendations, "")

	if strings.Contains(allRecommendations, ".NET") || strings.Contains(allRecommendations, "C#") {
		data.makeRecommendation("Review the .NET packaging cheatsheet: https://nhinv11.github.io/#/.NET")
		data.makeRecommendation("Consider using the unofficial JavaScript/TypeScript packaging tool: https://github.com/nhinv11/veracode-dotnet-packager")
	}

	if strings.Contains(allRecommendations, "Java ") {
		data.makeRecommendation("Review the Java packaging cheatsheet: https://nhinv11.github.io/#/Java")
	}

	data.makeRecommendation(fmt.Sprintf("Review the module configuration: %s", data.DetailedReport.getReviewModulesUrl(region)))
	data.makeRecommendation("Read the module selection guidance: https://community.veracode.com/s/article/What-are-Modules-and-how-do-my-results-change-based-on-what-I-select")

	data.makeRecommendation("Consider scheduling a consultation to review the packaging: https://docs.veracode.com/r/t_schedule_consultation")

	var report strings.Builder
	for _, recommendation := range *data.Recommendations {
		report.WriteString(formatRecommendationStringFormat(
			"%s\n",
			recommendation))
	}

	printTitle("Recommendations")
	colorPrintf(report.String() + "\n")
}
