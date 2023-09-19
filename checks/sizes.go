package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
)

func sizes(r *report.Report) {
	totalModuleSize(r)
	analysisSize(r)
}

func totalModuleSize(r *report.Report) {
	var totalSize = 0

	for _, module := range r.Modules {
		for _, instance := range module.Instances {
			totalSize += instance.SizeBytes
		}
	}

	if totalSize <= utils.MaximumTotalModuleSizeBytesThreshold {
		return
	}

	r.ReportIssue(fmt.Sprintf("The total size of all the modules was %s. This is a very large size and will likely take a long time to upload and scan.", utils.FormatBytes(uint64(totalSize))), report.IssueSeverityMedium)
	r.MakeRecommendation("Ensure you are not uploading more files than can be analysed by Veracode SAST.")
	r.MakeRecommendation("Follow the packaging guidance for each supported technology present within the application, as documented here: https://docs.veracode.com/r/compilation_packaging. Note there is also a useful cheat sheet which provides bespoke recommendations based off some questions about the application: https://docs.veracode.com/cheatsheet/.")
}

func analysisSize(r *report.Report) {
	if r.Scan.AnalysisSize <= utils.MaximumAnalysisSizeBytesThreshold {
		return
	}

	r.ReportIssue(fmt.Sprintf("The analysis size of the scan was %s. This is a very large size and will likely take a long time to upload and scan. Check that you are not selecting too many components for analysis.", utils.FormatBytes(r.Scan.AnalysisSize)), report.IssueSeverityMedium)
	r.MakeRecommendation("Ensure the correct modules have been selected for analysis and that the packaging guidance has been followed.")
	r.MakeRecommendation("Follow the packaging guidance for each supported technology present within the application, as documented here: https://docs.veracode.com/r/compilation_packaging. Note there is also a useful cheat sheet which provides bespoke recommendations based off some questions about the application: https://docs.veracode.com/cheatsheet/.")
}
