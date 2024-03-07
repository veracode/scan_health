package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veracode/scan_health/v2/report"
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

		assert.Empty(t, mockReport.Issues)
	})

	t.Run("Previous Scan doesn't exist", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: true}}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 333333, Name: "file3.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		var mockPreviousReport report.Report // Uninitialised
		previousScan(&mockReport, &mockPreviousReport)

		assert.Empty(t, mockReport.Issues)
	})

	t.Run("Module Count Has Changed", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Scan: report.Scan{
				BuildId: 1,
			},
			Modules: []report.Module{
				{Name: "file3.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: true}}},
				{Name: "file4.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: false}}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 333333, Name: "file3.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 444444, Name: "file4.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		mockPreviousReport := report.Report{
			Scan: report.Scan{
				BuildId: 2,
			},
			Modules: []report.Module{
				{Name: "file3.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: true}}},
				{Name: "file4.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: false}}},
				{Name: "file5.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: false}}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 333333, Name: "file3.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 444444, Name: "file4.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 555555, Name: "file5.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		previousScan(&mockReport, &mockPreviousReport)

		assert.Len(t, mockReport.Issues, 1)
		assert.Len(t, mockReport.Recommendations, 2)
	})

	t.Run("Module Selection Count Has Changed", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Scan: report.Scan{
				BuildId: 1,
			},
			Modules: []report.Module{
				{Name: "file3.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: true, Source: report.DetailedReportModuleSelected}}},
				{Name: "file4.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: false, Source: report.DetailedReportModuleSelected}}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 333333, Name: "file3.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 444444, Name: "file4.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		mockPreviousReport := report.Report{
			Scan: report.Scan{
				BuildId: 2,
			},
			Modules: []report.Module{
				{Name: "file3.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: true}}},
				{Name: "file4.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: true}}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 333333, Name: "file3.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 444444, Name: "file4.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		previousScan(&mockReport, &mockPreviousReport)

		assert.Len(t, mockReport.Issues, 1)
		assert.Len(t, mockReport.Recommendations, 2)
	})

	t.Run("Module Selection and Module Count Has Changed", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Scan: report.Scan{
				BuildId: 1,
			},
			Modules: []report.Module{
				{Name: "file3.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: true, Source: report.DetailedReportModuleSelected}}},
				{Name: "file4.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: false, Source: report.DetailedReportModuleSelected}}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 333333, Name: "file3.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 444444, Name: "file4.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		mockPreviousReport := report.Report{
			Scan: report.Scan{
				BuildId: 2,
			},
			Modules: []report.Module{
				{Name: "file3.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: true}}},
				{Name: "file4.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: false}}},
				{Name: "file5.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: true}}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 333333, Name: "file3.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 444444, Name: "file4.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 555555, Name: "file5.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		previousScan(&mockReport, &mockPreviousReport)

		assert.Len(t, mockReport.Issues, 2)
		assert.Len(t, mockReport.Recommendations, 3)
	})

	/*
			t.Run("Module Selection Has Changed (Simple, no versions)", func(t *testing.T) {
				t.Parallel()
				mockReport := report.Report{
					Scan: report.Scan{
						BuildId: 1,
					},
					Modules: []report.Module{
						{Name: "file_three.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: true, Id: 100}}},
						{Name: "file_four.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: false, Id: 101}}},
					},
					UploadedFiles: []report.UploadedFile{
						{Id: 333333, Name: "file_three.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
						{Id: 444444, Name: "file_four.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
					},
				}

				mockPreviousReport := report.Report{
					Scan: report.Scan{
						BuildId: 2,
					},
					Modules: []report.Module{
						{Name: "file_three.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: true, Id: 200}}},
						{Name: "file_four.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: true, Id: 201}}},
					},
					UploadedFiles: []report.UploadedFile{
						{Id: 333333, Name: "file_three.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
						{Id: 444444, Name: "file_four.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
					},
				}

				previousScan(&mockReport, &mockPreviousReport)

				assert.Len(t, mockReport.Issues, 1)
				assert.Len(t, mockReport.Recommendations, 1)
			})

		t.Run("Module Selection Has Not Changed (Modules have Versions)", func(t *testing.T) {
			t.Parallel()
			mockReport := report.Report{
				Scan: report.Scan{
					BuildId: 1,
				},
				Modules: []report.Module{
					{Name: "file3-12.3.4.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: true, Id: 100}}},
					{Name: "file4-3.4.5.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: false, Id: 101}}},
				},
				UploadedFiles: []report.UploadedFile{
					{Id: 333333, Name: "file3.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
					{Id: 444444, Name: "file4.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				},
			}

			mockPreviousReport := report.Report{
				Scan: report.Scan{
					BuildId: 2,
				},
				Modules: []report.Module{
					{Name: "file3-12.3.5.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: true, Id: 200}}},
					{Name: "file4-3.4.6.jar", Instances: []report.ModuleInstance{{IsDependency: false, IsSelected: false, Id: 201}}},
				},
				UploadedFiles: []report.UploadedFile{
					{Id: 333333, Name: "file3.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
					{Id: 444444, Name: "file4.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				},
			}

			previousScan(&mockReport, &mockPreviousReport)

			assert.Empty(t, mockReport.Issues)
			assert.Empty(t, mockReport.Recommendations)
		})
	*/
	t.Run("Filename normalization on multiple versions", func(t *testing.T) {
		t.Parallel()

		filenames := []string{
			"file-1.0.0.jar",
			"file-1.0.1-SNAPSHOT.jar",
			"file-2.0.jar",
			"file.jar",
		}

		normalizedFiles := make(map[string]int)

		for _, filename := range filenames {
			normalized := normalizeFilename(filename)
			normalizedFiles[normalized]++
		}

		assert.Equal(t, len(normalizedFiles), 1)
		assert.Equal(t, normalizedFiles["file.jar"], 4)
	})

	t.Run("Filename normalization on files that include a number", func(t *testing.T) {
		t.Parallel()

		filenames := []string{
			"file2test-1.0.0.dll",
			"file2test-2.0.0.dll",
		}

		normalizedFiles := make(map[string]int)

		for _, filename := range filenames {
			normalized := normalizeFilename(filename)
			normalizedFiles[normalized]++
		}

		assert.Equal(t, len(normalizedFiles), 1)
		assert.Equal(t, normalizedFiles["file2test.dll"], 2)
	})
}
