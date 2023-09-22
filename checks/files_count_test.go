package checks

import (
	"github.com/antfie/scan_health/v2/utils"
	"testing"

	"github.com/antfie/scan_health/v2/report"
	"github.com/stretchr/testify/assert"
)

func createUploadedFiles(numFiles int) []report.UploadedFile {
	var files []report.UploadedFile

	for i := 0; i < numFiles; i++ {
		files = append(files, report.UploadedFile{Id: i})
	}

	return files
}

// Test Cases
func TestTooManyModules(t *testing.T) {

	// Test Case 1: No duplicates
	t.Run("Typical number of modules", func(t *testing.T) {
		t.Parallel()

		testReport := report.Report{
			UploadedFiles: createUploadedFiles(300),
			Issues:        []report.Issue{},
		}

		tooManyFilesUploaded(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	t.Run("Typical number of modules", func(t *testing.T) {
		t.Parallel()

		testReport := report.Report{
			UploadedFiles: createUploadedFiles(utils.MaximumUploadedFileCountThreshold + 1),
			Issues:        []report.Issue{},
		}

		tooManyFilesUploaded(&testReport)
		assert.Equal(t, 1, len(testReport.Issues))
		assert.Equal(t, report.IssueSeverityMedium, testReport.Issues[0].Severity)
	})
}
