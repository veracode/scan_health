package checks

import (
	"testing"

	"github.com/antfie/scan_health/v2/report"
	"github.com/stretchr/testify/assert"
)

// Test Cases
func TestUnexpectedSourceCode(t *testing.T) {

	// Test Case 1: No Roslyn Files
	t.Run("No Unexpected Files", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		unexpectedSourceCode(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	t.Run("4 files uploaded, 1 unexpected Files", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.jar", MD5: "hash1", IsIgnored: false},
				{Id: 222222, Name: "file2.cs", MD5: "hash2", IsIgnored: false},
				{Id: 333333, Name: "file3.dll", MD5: "hash3", IsIgnored: false},
				{Id: 444444, Name: "file4.exe", MD5: "hash4", IsIgnored: false},
			},
			Issues: []report.Issue{},
		}

		unexpectedSourceCode(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Issues[0].Description, "A C# source code file")
		assert.Equal(t, 2, len(testReport.Recommendations))
	})

	t.Run("4 files uploaded, 2 unexpected Files (same type)", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.jar", MD5: "hash1", IsIgnored: false},
				{Id: 222222, Name: "file2.cs", MD5: "hash2", IsIgnored: false},
				{Id: 333333, Name: "file3.cs", MD5: "hash3", IsIgnored: false},
				{Id: 444444, Name: "file4.exe", MD5: "hash4", IsIgnored: false},
			},
			Issues: []report.Issue{},
		}

		unexpectedSourceCode(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Issues[0].Description, "2 C# source code files")
		assert.Equal(t, 2, len(testReport.Recommendations))
	})

	t.Run("4 files uploaded, 2 unexpected Files (different types)", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.jar", MD5: "hash1", IsIgnored: false},
				{Id: 222222, Name: "file2.cs", MD5: "hash2", IsIgnored: false},
				{Id: 333333, Name: "file3.java", MD5: "hash3", IsIgnored: false},
				{Id: 444444, Name: "file4.exe", MD5: "hash4", IsIgnored: false},
			},
			Issues: []report.Issue{},
		}

		unexpectedSourceCode(&testReport)
		if !assert.Equal(t, 2, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Issues[0].Description, "A Java source code file")
		assert.Contains(t, testReport.Issues[1].Description, "A C# source code file")
		if !assert.Equal(t, 4, len(testReport.Recommendations)) {
			t.FailNow()
		}
	})
}
