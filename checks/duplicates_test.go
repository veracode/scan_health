package checks

import (
	"testing"

	"github.com/antfie/scan_health/v2/report"
	"github.com/stretchr/testify/assert"
)

// Test Cases
func TestIdenticalModulesParallel(t *testing.T) {

	// Test Case 1: No duplicates
	t.Run("No duplicates", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		duplicateModules(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	// Test Case 2: Same filename, same hash
	t.Run("Files with same name same hash", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file1", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
			},
		}

		duplicateModules(&testReport)
		assert.Equal(t, 1, len(testReport.Issues))
		assert.Equal(t, report.IssueSeverityMedium, testReport.Issues[0].Severity)
	})

	// Test Case 3: Same filename, different hashes
	t.Run("Files with same name different hashes", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file1", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		duplicateModules(&testReport)
		assert.Equal(t, 1, len(testReport.Issues))
		assert.Equal(t, report.IssueSeverityHigh, testReport.Issues[0].Severity)
	})
}
