package checks

import (
	"strings"
	"testing"

	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
	"github.com/stretchr/testify/assert"
)

// Test Cases
func TestSizes(t *testing.T) {

	// Test Case 1: No duplicates
	t.Run("Normal Upload", func(t *testing.T) {
		t.Parallel()

		testReport := report.Report{
			Scan: report.Scan{
				AnalysisSize: 10 * 1024 * 1024,
			},
			Modules: []report.Module{
				{Name: "module1",
					IsThirdParty: false,
					IsIgnored:    false,
					Instances: []report.ModuleInstance{{
						SizeBytes:      1024 * 1024,
						IsSelected:     true,
						HasFatalErrors: false,
					}},
				},
				{Name: "module2",
					IsThirdParty: false,
					IsIgnored:    false,
					Instances: []report.ModuleInstance{{
						SizeBytes:      2048 * 1024,
						IsSelected:     true,
						HasFatalErrors: false,
					}},
				},
			},
			Issues: []report.Issue{},
		}

		sizes(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	t.Run("Overlarge Analysis Size", func(t *testing.T) {
		t.Parallel()

		testReport := report.Report{
			Scan: report.Scan{
				AnalysisSize: utils.MaximumAnalysisSizeBytesThreshold + 1,
			},
			Modules: []report.Module{},
			Issues:  []report.Issue{},
		}

		sizes(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.True(t, strings.Contains(testReport.Issues[0].Description, "very large"))

		// We do actually get 2 recommendations
		assert.Equal(t, 2, len(testReport.Recommendations))
	})

	t.Run("Cumulatively Large Scan", func(t *testing.T) {
		t.Parallel()

		individualModuleSize := (utils.MaximumTotalModuleSizeBytesThreshold / 3) + 1
		testReport := report.Report{
			Scan: report.Scan{},
			Modules: []report.Module{
				{Name: "module1", IsIgnored: false, IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{
							Id:        1,
							SizeBytes: individualModuleSize,
						},
					},
				},
				{Name: "module1", IsIgnored: false, IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{
							Id:        3,
							SizeBytes: individualModuleSize,
						},
					},
				},
				{Name: "module1", IsIgnored: false, IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{
							Id:        3,
							SizeBytes: individualModuleSize,
						},
					},
				},
			},
			Issues: []report.Issue{},
		}

		sizes(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.True(t, strings.Contains(testReport.Issues[0].Description, "The total size of all the modules was"))

		// We do actually get 2 recommendations
		assert.Equal(t, 2, len(testReport.Recommendations))
	})
}
