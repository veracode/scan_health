package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veracode/scan_health/v2/report"
)

func TestUnselectedFirstPartyFiles(t *testing.T) {

	t.Run("1st party top level module selected", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.exe",
					IsIgnored:    false,
					IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: true},
					}},
			},
		}

		unselectedFirstParty(&mockReport)
		assert.Empty(t, mockReport.Issues)
	})

	t.Run("3rd party top level module not selected", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.exe",
					IsIgnored:    false,
					IsThirdParty: true,
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: false},
					}},
			},
		}

		unselectedFirstParty(&mockReport)
		assert.Empty(t, mockReport.Issues)
	})

	t.Run("1st party dependency not selected", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.exe",
					IsIgnored:    false,
					IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{IsDependency: true},
						{IsSelected: false},
					}},
			},
		}

		unselectedFirstParty(&mockReport)
		assert.Empty(t, mockReport.Issues)
	})

	t.Run("1st party top level module not selected", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.exe",
					IsIgnored:    false,
					IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: false},
					}},
			},
		}

		unselectedFirstParty(&mockReport)

		if !assert.Equal(t, 1, len(mockReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, mockReport.Issues[0].Description, "A potential first-party module")
		assert.Equal(t, 1, len(mockReport.Recommendations))
	})

	t.Run("2 1st party top level modules not selected", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.exe",
					IsIgnored:    false,
					IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: false},
					}},
				{Name: "file4.dll",
					IsIgnored:    false,
					IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: false},
					}},
			},
		}

		unselectedFirstParty(&mockReport)

		if !assert.Equal(t, 1, len(mockReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, mockReport.Issues[0].Description, "2 potential first-party modules")
		assert.Equal(t, 1, len(mockReport.Recommendations))
	})
}
