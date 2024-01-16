package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veracode/scan_health/v2/report"
)

// Test Cases
func TestUnwantedFiles(t *testing.T) {

	t.Run("No sensitive files", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		detectUnwantedFiles(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	t.Run("Many Files, A Single Unwanted File", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "test.7z", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		detectUnwantedFiles(&testReport)

		assert.Equal(t, 1, len(testReport.Issues))
		assert.Contains(t, testReport.Issues[0].Description, "A 7-zip file was uploaded: \"test.7z\"")
		assert.Equal(t, 2, len(testReport.Recommendations))
	})

	t.Run("Many Files, Many Unwanted Files", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "test.7z", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "test.pyc", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		detectUnwantedFiles(&testReport)

		assert.Equal(t, 2, len(testReport.Issues))
		assert.Contains(t, testReport.Issues[0].Description, "A 7-zip file was uploaded: \"test.7z\"")
		assert.Contains(t, testReport.Issues[1].Description, "A compiled Python file was uploaded: \"test.pyc\"")
		assert.Equal(t, 3, len(testReport.Recommendations))
	})
}
