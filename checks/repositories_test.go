package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veracode/scan_health/v2/report"
)

// Test Cases
func TestUploadedRepositoryFiles(t *testing.T) {

	// Test Case 1: No duplicates
	t.Run("No repository files", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		repositories(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	// Test Case 2: Same filename, same hash
	t.Run("Watchman Canary", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "fsmonitor-watchman.sample", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
			},
		}

		repositories(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}
		assert.Equal(t, report.IssueSeverityMedium, testReport.Issues[0].Severity)
		assert.Equal(t, 1, len(testReport.Recommendations))
	})

	// Test Case 3: Same filename, different hashes
	t.Run("FETCH_HEAD canary", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "FETCH_HEAD", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		repositories(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}
		assert.Equal(t, report.IssueSeverityMedium, testReport.Issues[0].Severity)
		assert.Equal(t, 1, len(testReport.Recommendations))
	})
}
