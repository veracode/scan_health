package checks

import (
	"github.com/antfie/scan_health/v2/utils"
	"strings"
	"testing"

	"github.com/antfie/scan_health/v2/report"
	"github.com/stretchr/testify/assert"
)

// Test Cases
func TestTooManyFlaws(t *testing.T) {

	// Test Case 1: No duplicates
	t.Run("Reasonable number of flaws", func(t *testing.T) {
		t.Parallel()

		testReport := report.Report{
			Flaws: report.FlawSummary{
				Total: 200,
			},
			Issues: []report.Issue{},
		}

		flawCount(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	t.Run("Zero flaws", func(t *testing.T) {
		t.Parallel()

		testReport := report.Report{
			Flaws: report.FlawSummary{
				Total: 0,
			},
			Issues: []report.Issue{},
		}

		flawCount(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.True(t, strings.Contains(testReport.Issues[0].Description, "No flaws"))

		if !assert.Equal(t, 1, len(testReport.Recommendations)) {
			t.FailNow()
		}

		assert.True(t, strings.Contains(testReport.Recommendations[0], "When no flaws have been found"))
	})

	t.Run("Too Many Flaws", func(t *testing.T) {
		t.Parallel()

		testReport := report.Report{
			Flaws: report.FlawSummary{
				Total: utils.MaximumFlawCountThreshold + 1,
			},
			Issues: []report.Issue{},
		}

		flawCount(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.True(t, strings.Contains(testReport.Issues[0].Description, "A large number"))

		if !assert.Equal(t, 1, len(testReport.Recommendations)) {
			t.FailNow()
		}

		assert.True(t, strings.Contains(testReport.Recommendations[0], "scan could be misconfigured"))
	})
}
