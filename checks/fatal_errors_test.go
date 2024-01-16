package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veracode/scan_health/v2/report"
)

// Test Cases
func TestModulesAreFine(t *testing.T) {

	// Test Case 1: No duplicates
	t.Run("No Module Errors", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "module1",
					IsThirdParty: false,
					IsIgnored:    false,
					Instances: []report.ModuleInstance{{
						IsSelected:     true,
						HasFatalErrors: false,
					}},
				},
				{Name: "module2",
					IsThirdParty: false,
					IsIgnored:    false,
					Instances: []report.ModuleInstance{{
						IsSelected:     true,
						HasFatalErrors: false,
					}},
				},
			},
			Issues: []report.Issue{},
		}

		fatalErrors(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	t.Run("Two Windows Modules With Missing Primary Debug Symbols", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "module1.exe",
					IsThirdParty: false,
					IsIgnored:    false,
					Instances: []report.ModuleInstance{{
						IsSelected:     true,
						HasFatalErrors: true,
						Status:         "Primary Files Compiled without Debug Symbols",
					}},
				},
				{Name: "module2.dll",
					IsThirdParty: false,
					IsIgnored:    false,
					Instances: []report.ModuleInstance{{
						IsSelected:     true,
						HasFatalErrors: true,
						Status:         "Primary Files Compiled without Debug Symbols",
					}},
				},
				{Name: "module3",
					IsThirdParty: false,
					IsIgnored:    false,
					Instances: []report.ModuleInstance{{
						IsSelected:     true,
						HasFatalErrors: false,
					}},
				},
			},
			Issues: []report.Issue{},
		}

		fatalErrors(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		if !assert.Contains(t, testReport.Issues[0].Description, "2 modules") {
			t.FailNow()
		}

		assert.Equal(t, report.IssueSeverityHigh, testReport.Issues[0].Severity)

		if !assert.Equal(t, 1, len(testReport.Recommendations)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Recommendations[0], "PDB")
	})

	t.Run("Two Java Modules With No Scannable Binaries", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "module1.war",
					IsThirdParty: false,
					IsIgnored:    false,
					Instances: []report.ModuleInstance{{
						IsSelected:     false,
						HasFatalErrors: true,
						Status:         "No Scannable Binaries",
					}},
				},
				{Name: "module2.ear",
					IsThirdParty: false,
					IsIgnored:    false,
					Instances: []report.ModuleInstance{{
						IsSelected:     false,
						HasFatalErrors: true,
						Status:         "No Scannable Binaries",
					}},
				},
				{Name: "module3.jar",
					IsThirdParty: false,
					IsIgnored:    false,
					Instances: []report.ModuleInstance{{
						IsSelected:     true,
						HasFatalErrors: true,
						Status:         "No Scannable Binaries",
					}},
				},
			},
			Issues: []report.Issue{},
		}

		fatalErrors(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		if !assert.Contains(t, testReport.Issues[0].Description, "3 Java modules") {
			t.FailNow()
		}

		assert.Equal(t, report.IssueSeverityHigh, testReport.Issues[0].Severity)

		if !assert.Equal(t, 1, len(testReport.Recommendations)) {
			t.FailNow()
		}
	})

	t.Run("Two Java Modules With No Scannable Binaries", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "module1.war",
					IsThirdParty: false,
					IsIgnored:    false,
					Instances: []report.ModuleInstance{{
						IsSelected:     false,
						HasFatalErrors: true,
						Status:         "does not support jar files nested inside",
					}},
				},
			},
			Issues: []report.Issue{},
		}

		fatalErrors(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		if !assert.Contains(t, testReport.Issues[0].Description, "nested/shaded") {
			t.FailNow()
		}

		assert.Equal(t, report.IssueSeverityHigh, testReport.Issues[0].Severity)

		if !assert.Equal(t, 1, len(testReport.Recommendations)) {
			t.FailNow()
		}
	})
}
