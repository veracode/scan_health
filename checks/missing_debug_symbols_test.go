package checks

import (
	"testing"

	"github.com/antfie/scan_health/v2/report"
	"github.com/stretchr/testify/assert"
)

func TestMissingDebugSymbols(t *testing.T) {

	t.Run("No Executables/DLLs", func(t *testing.T) {
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

		missingDebugSymbols(&mockReport)

		assert.Empty(t, mockReport.Issues)
	})

	t.Run("Executables/DLLs but no issues", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.dll",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
					}},
				{Name: "file4.exe",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
					}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file3.dll", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file4.exe", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		missingDebugSymbols(&mockReport)

		assert.Empty(t, mockReport.Issues)
	})

	t.Run("Java with missing debug", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.jar",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{Issues: []string{"No supporting files or PDB files"}},
					}},
				{Name: "file4.war",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{Issues: []string{"No supporting files or PDB files"}},
					}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file3.exe", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file4.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 111111, Name: "file3.jar", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file4.war", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		missingDebugSymbols(&mockReport)

		assert.Empty(t, mockReport.Issues)
		assert.Empty(t, mockReport.Recommendations)
	})

	t.Run("Executables + DLLs with missing debug", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.exe",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{Issues: []string{"No supporting files or PDB files"}},
					}},
				{Name: "file4.dll",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{Issues: []string{"No supporting files or PDB files"}},
					}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file3.exe", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file4.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		missingDebugSymbols(&mockReport)

		if !assert.Equal(t, 1, len(mockReport.Issues)) {
			assert.Contains(t, mockReport.Issues[0].Description, "2 modules")
		}

		assert.Equal(t, 1, len(mockReport.Recommendations))
	})
}
