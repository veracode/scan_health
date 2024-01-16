package data

import (
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
)

func (api API) GetPreviousBuildId(r *report.Report) int {
	for _, build := range r.OtherScans {
		buildInfo := api.GetBuildInfo(r.Scan.ApplicationId, build.BuildId, build.SandboxId)

		if buildInfo.Build.AnalysisUnit.Status != "Results Ready" {
			continue
		}

		if len(buildInfo.Build.AnalysisUnit.PublishedDate) == 0 {
			continue
		}

		publishedDate := utils.ParseVeracodeDate(buildInfo.Build.AnalysisUnit.PublishedDate).Local()
		if r.Scan.PublishedDate.After(publishedDate) {
			return build.BuildId
		}
	}

	return 0
}
