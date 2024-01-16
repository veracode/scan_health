package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veracode/scan_health/v2/report"
)

func createReportWithSingleSensitiveFile(extension string) report.Report {
	return report.Report{
		UploadedFiles: []report.UploadedFile{
			{Id: 111111, Name: "file1.exe", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
			{Id: 222222, Name: "file2" + extension, MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			{Id: 333333, Name: "file3.jar", MD5: "hash3", IsIgnored: false, IsThirdParty: false},
		},
	}
}

func testAndAssertByExtension(t *testing.T, extension string) {
	testReport := createReportWithSingleSensitiveFile(extension)

	sensitiveFiles(&testReport)
	if !assert.Equal(t, 1, len(testReport.Issues)) {
		t.FailNow()
	}
	assert.Equal(t, report.IssueSeverityHigh, testReport.Issues[0].Severity)

	assert.Equal(t, 2, len(testReport.Recommendations))
}

// Test Cases
func TestSensitiveFilesParallel(t *testing.T) {

	t.Run("No sensitive files", func(t *testing.T) {
		t.Parallel()
		testReport := createReportWithSingleSensitiveFile("")

		duplicateModules(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	t.Run("Many Files, A Single Secret File", func(t *testing.T) {
		t.Parallel()

		testAndAssertByExtension(t, ".p12")
	})

	t.Run("Many Files, A Single Backup File", func(t *testing.T) {
		t.Parallel()

		testAndAssertByExtension(t, ".old")
	})

	t.Run("Many Files, A Single Word File", func(t *testing.T) {
		t.Parallel()

		testAndAssertByExtension(t, ".docm")
	})

	t.Run("Many Files, A Single Excel File", func(t *testing.T) {
		t.Parallel()

		testAndAssertByExtension(t, ".xlsx")
	})

	t.Run("Many Files, A Single Jupyter Notebook File", func(t *testing.T) {
		t.Parallel()

		testAndAssertByExtension(t, ".ipynb")
	})

	t.Run("Many Sensitive Files", func(t *testing.T) {
		t.Parallel()

		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.p12", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2.old", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 333333, Name: "file3.docm", MD5: "hash3", IsIgnored: false, IsThirdParty: false},
				{Id: 444444, Name: "file3.xlsx", MD5: "hash4", IsIgnored: false, IsThirdParty: false},
				{Id: 555555, Name: "file3.ipynb", MD5: "hash5", IsIgnored: false, IsThirdParty: false},
			},
		}

		sensitiveFiles(&testReport)
		if !assert.Equal(t, 5, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Equal(t, report.IssueSeverityHigh, testReport.Issues[0].Severity)

		assert.Equal(t, 5, len(testReport.Recommendations))
	})
}
