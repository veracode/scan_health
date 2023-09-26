package checks

import (
	"testing"

	"github.com/antfie/scan_health/v2/report"
	"github.com/stretchr/testify/assert"
)

func TestThirdParty(t *testing.T) {

	t.Run("No Third Party", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.jar",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
					}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 333333, Name: "file3.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		thirdParty(&mockReport)

		assert.Equal(t, len(mockReport.Issues), 0)
	})

	t.Run("Third Party uploaded, not selected", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.dll",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
					}},
				{Name: "Antlr3.Runtime.dll",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: false},
					}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "Antlr3.Runtime.dll", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file4.exe", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		thirdParty(&mockReport)

		assert.Equal(t, len(mockReport.Issues), 0)
	})

	t.Run("Third Party uploaded, selected", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.dll",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
					}},
				{Name: "Antlr3.Runtime.dll",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: true},
					}},
				{Name: "DevExpress.Components.dll",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: true},
					}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "Antlr3.Runtime.dll", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file4.exe", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 333333, Name: "DevExpress.Components.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		thirdParty(&mockReport)

		if !assert.Equal(t, len(mockReport.Issues), 1) {
			t.FailNow()
		}

		assert.Contains(t, mockReport.Issues[0].Description, "2 third-party components")

		assert.Equal(t, report.IssueSeverityMedium, mockReport.Issues[0].Severity)
		assert.Equal(t, len(mockReport.Recommendations), 1)
		assert.True(t, mockReport.UploadedFiles[0].IsThirdParty)
		assert.True(t, mockReport.UploadedFiles[2].IsThirdParty)
	})
}
