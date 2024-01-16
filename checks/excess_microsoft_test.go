package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veracode/scan_health/v2/report"
)

// Test Cases
func TestMicrosoftRuntimeFiles(t *testing.T) {

	// Test Case 1: No Roslyn Files
	t.Run("No Roslyn Files", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		excessMicrosoft(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	// Test Case 2: Our 'canary' Roslyn file (csc.exe) is present
	t.Run("Canary Roslyn Files", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "csc.exe", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file1", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		excessMicrosoft(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues), "Expected csc.exe in the file list, wasn't found") {
			t.FailNow()
		}

		assert.Equal(t, report.IssueSeverityMedium, testReport.Issues[0].Severity)
	})
}
