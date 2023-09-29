package checks

import (
	"testing"

	"github.com/antfie/scan_health/v2/report"
	"github.com/stretchr/testify/assert"
)

// Test Cases
func TestMissingSCAComponents(t *testing.T) {

	t.Run("Valid Files", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.dll", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Scan: report.Scan{
				IsSCADataAvailable: true,
			},
			SCAComponents: []string{"test1.dll", "file2.jar"},
			Issues:        []report.Issue{},
		}

		missingSCAComponents(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	t.Run("SCA Data is not available", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.dll", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Scan: report.Scan{
				IsSCADataAvailable: false,
			},
			Issues: []report.Issue{},
		}

		missingSCAComponents(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	t.Run("No SCA files were found", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "csc.exe", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file1", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
			},
			Scan: report.Scan{
				IsSCADataAvailable: true,
			},
			SCAComponents: []string{},
			Issues:        []report.Issue{},
		}

		missingSCAComponents(&testReport)
		assert.Equal(t, 1, len(testReport.Issues))
		assert.Equal(t, 1, len(testReport.Recommendations))
	})

}
