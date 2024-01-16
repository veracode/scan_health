package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veracode/scan_health/v2/report"
)

// Test Cases
func TestGeneralRecommendations(t *testing.T) {

	t.Run("No Recommendations Exists", func(t *testing.T) {
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

			Issues:          []report.Issue{},
			Recommendations: []string{},
		}

		generalRecommendations(&testReport)
		assert.Empty(t, testReport.Issues)
		assert.Empty(t, testReport.Recommendations)
	})

	t.Run("No Module Recommendation Exists", func(t *testing.T) {
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
			Recommendations: []string{
				"Test Recommendation",
			},
		}

		generalRecommendations(&testReport)
		assert.Empty(t, testReport.Issues)

		// No additional recommendations made
		assert.Equal(t, 1, len(testReport.Recommendations))
	})

	t.Run("Module Recommendation Exists", func(t *testing.T) {
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
			Recommendations: []string{
				"module",
			},
		}

		generalRecommendations(&testReport)
		assert.Empty(t, testReport.Issues)

		// Additional recommendations are made on the presence of 'module' in the recommendations list
		assert.Equal(t, 6, len(testReport.Recommendations))
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
					}},
				},
				{Name: "gradle-wrapper.jar",
					IsThirdParty: false,
					IsIgnored:    false,
					Instances: []report.ModuleInstance{{
						IsSelected:     true,
						HasFatalErrors: false,
					}},
				},
				{Name: "test.jar",
					IsThirdParty: false,
					IsIgnored:    false,
					Instances: []report.ModuleInstance{{
						IsSelected:     false,
						HasFatalErrors: false,
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

		assert.Equal(t, 1, len(testReport.Recommendations))
	})
}
