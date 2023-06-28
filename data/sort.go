package data

import (
	"github.com/antfie/scan_health/v2/report"
	"sort"
)

func sortData(r *report.Report) {
	// Sort uploaded files by name for consistency
	sort.Slice(r.UploadedFiles, func(i, j int) bool {
		return r.UploadedFiles[i].Name < r.UploadedFiles[j].Name
	})

	// Sort modules by name for consistency
	sort.Slice(r.Modules, func(i, j int) bool {
		return r.Modules[i].Name < r.Modules[j].Name
	})
}
