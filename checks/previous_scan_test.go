package checks

import (
	"testing"

	"github.com/antfie/scan_health/v2/report"
	"github.com/stretchr/testify/assert"
)

func TestPreviousScan(t *testing.T) {

	t.Run("Previous Scan is Identical", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: true}}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 333333, Name: "file3.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		mockPreviousReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: true}}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 333333, Name: "file3.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		previousScan(&mockReport, &mockPreviousReport)

		assert.Equal(t, len(mockReport.Issues), 0)
	})
}
