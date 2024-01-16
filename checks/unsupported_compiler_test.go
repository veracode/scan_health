package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veracode/scan_health/v2/report"
)

func TestUnsupportedCompiler(t *testing.T) {

	t.Run("All Supported Modules", func(t *testing.T) {
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

		unsupportedPlatformOrCompiler(&mockReport)
		assert.Empty(t, mockReport.Issues)
	})

	t.Run("Ignores third party files", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.exe",
					IsIgnored:    false,
					IsThirdParty: true,
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: true},
					}},
			},
		}

		unsupportedPlatformOrCompiler(&mockReport)
		assert.Empty(t, mockReport.Issues)
	})

	t.Run("One unsupported module", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "bob.exe",
					IsIgnored:    false,
					IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: true},
						{HasFatalErrors: true},
						{Status: "(Fatal)Unsupported Platform"},
					}},
			},
		}

		unsupportedPlatformOrCompiler(&mockReport)
		if !assert.Equal(t, 1, len(mockReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, mockReport.Issues[0].Description, "A module could not be scanned")
		assert.Equal(t, 1, len(mockReport.Recommendations))
	})

	t.Run("2 unsupported modules", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "bob.exe",
					IsIgnored:    false,
					IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: true},
						{HasFatalErrors: true},
						{Status: "(Fatal)Unsupported Platform"},
					}},

				{Name: "bob2.exe",
					IsIgnored:    false,
					IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: true},
						{HasFatalErrors: true},
						{Status: "(Fatal)Unsupported Compiler"},
					}},
			},
		}

		unsupportedPlatformOrCompiler(&mockReport)
		if !assert.Equal(t, 1, len(mockReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, mockReport.Issues[0].Description, "2 modules could not be scanned")
		assert.Equal(t, 1, len(mockReport.Recommendations))
	})

}
