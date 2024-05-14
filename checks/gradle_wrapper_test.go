package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veracode/scan_health/v2/report"
)

// Test Cases
func TestGradleWrapper(t *testing.T) {

	// Test Case 1: No Roslyn Files
	t.Run("No Gradle Wrapper Module", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "test.jar",
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

		gradleWrapper(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	t.Run("Gradle Wrapper is the only selected one", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "test.jar",
					IsThirdParty: false,
					IsIgnored:    false,
					Instances: []report.ModuleInstance{{
						IsSelected:     false,
						HasFatalErrors: false,
						Source:         report.DetailedReportModuleSelected,
					}},
				},
				{Name: "gradle-wrapper.jar",
					IsThirdParty: false,
					IsIgnored:    false,
					Instances: []report.ModuleInstance{{
						IsSelected:     true,
						HasFatalErrors: false,
						Source:         report.DetailedReportModuleSelected,
					}},
				},
				{Name: "test.jar",
					IsThirdParty: false,
					IsIgnored:    false,
					Instances: []report.ModuleInstance{{
						IsSelected:     false,
						HasFatalErrors: false,
						Source:         report.DetailedReportModuleSelected,
					}},
				},
			},
			Issues: []report.Issue{},
		}

		gradleWrapper(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Issues[0].Description, "The only module selected ")

		assert.Equal(t, 2, len(testReport.Recommendations))
	})
}
