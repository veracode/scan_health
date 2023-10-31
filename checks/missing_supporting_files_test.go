package checks

import (
	"testing"

	"github.com/antfie/scan_health/v2/report"
	"github.com/stretchr/testify/assert"
)

// Test Cases
func TestMissingSupportingFiles(t *testing.T) {

	t.Run("Valid Files and Modules", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "file1.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: true}}},
				{Name: "file2.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: true}}},
				{Name: "file3.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: false}}},
			},
			Issues: []report.Issue{},
		}

		missingSupportingFiles(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	t.Run("Module with 1 missing file", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "file1.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: true}}},
				{Name: "file2.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: true}, {Issues: []string{
					"Missing Supporting Files - 1 File",
				}}}},
				{Name: "file3.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: false}}},
			},
			Issues: []report.Issue{},
		}

		missingSupportingFiles(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Issues[0].Description, "A selected module \"file2.dll\" was found to be missing 1 file.")
		assert.Equal(t, 2, len(testReport.Recommendations))
	})

	t.Run("Module with 3 missing files", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "file1.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: true}}},
				{Name: "file2.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: true}, {Issues: []string{
					"Missing Supporting Files - 3 Files",
				}}}},
				{Name: "file3.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: false}}},
			},
			Issues: []report.Issue{},
		}

		missingSupportingFiles(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Issues[0].Description, "A selected module \"file2.dll\" was found to be missing 3 files.")
		assert.Equal(t, 2, len(testReport.Recommendations))
	})

	t.Run("Module with 4 'missing 1 file'", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "file1.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: true}}},
				{Name: "file2.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: true}, {Issues: []string{
					"Missing Supporting Files - 1 File",
					"Missing Supporting Files - 1 File",
					"Missing Supporting Files - 1 File",
					"Missing Supporting Files - 1 File",
				}}}},
				{Name: "file3.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: false}}},
			},
			Issues: []report.Issue{},
		}

		missingSupportingFiles(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Issues[0].Description, "A selected module \"file2.dll\" was found to be missing 4 files.")
		assert.Equal(t, 2, len(testReport.Recommendations))
	})

	t.Run("Multiple Modules missing 1 file", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "file1.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: true}, {Issues: []string{
					"Missing Supporting Files - 1 File",
				}}}},
				{Name: "file2.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: true}, {Issues: []string{
					"Missing Supporting Files - 1 File",
				}}}},
				{Name: "file3.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: false}}},
			},
			Issues: []report.Issue{},
		}

		missingSupportingFiles(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Issues[0].Description, "2 selected modules were selected as entry points that were found to be missing a total of 2 files")
		assert.Equal(t, 2, len(testReport.Recommendations))
	})
}
