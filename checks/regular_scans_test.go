package checks

import (
	"github.com/antfie/scan_health/v2/utils"
	"testing"
	"time"

	"github.com/antfie/scan_health/v2/report"
	"github.com/stretchr/testify/assert"
)

// Test Cases
func TestScanFrequency(t *testing.T) {

	// Test Case 1: No duplicates
	t.Run("Recent Scan", func(t *testing.T) {
		t.Parallel()

		testReport := report.Report{
			LastAppActivity: time.Now(),
			Issues:          []report.Issue{},
		}

		regularScans(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	t.Run("No scan", func(t *testing.T) {
		t.Parallel()

		testReport := report.Report{
			LastAppActivity: time.UnixMilli(0),
			Issues:          []report.Issue{},
		}

		regularScans(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	t.Run("Too Long Ago", func(t *testing.T) {
		t.Parallel()

		testReport := report.Report{
			LastAppActivity: time.Now().AddDate(0, 0, -(utils.NotUsingAutomationIfScanOlderThanDays + 1)),
			Issues:          []report.Issue{},
		}

		regularScans(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Issues[0].Description, "last scanned over")

		assert.Equal(t, 1, len(testReport.Recommendations))
	})
}
