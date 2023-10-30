package checks

import (
	"testing"

	"github.com/antfie/scan_health/v2/report"
	"github.com/stretchr/testify/assert"
)

func TestNestedArchives(t *testing.T) {

	t.Run("Normal Files", func(t *testing.T) {
		t.Parallel()

		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 1, Name: "test.jar"},
				{Id: 2, Name: "test2.jar"},
			},
			Issues: []report.Issue{},
		}

		nestedArchives(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	t.Run("Normal Files and a nested Archive", func(t *testing.T) {
		t.Parallel()

		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 1, Name: "test.jar"},
				{Id: 2, Name: "test2.jar"},
				{Id: 3, Name: "test3.jar", Status: "Archive File Within Another Archive"},
			},
			Issues: []report.Issue{},
		}

		nestedArchives(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Issues[0].Description, "A nested archive")

		assert.Equal(t, 1, len(testReport.Recommendations))
	})

	t.Run("Normal Files and multiple nested Archive", func(t *testing.T) {
		t.Parallel()

		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 1, Name: "test.jar"},
				{Id: 2, Name: "test2.jar"},
				{Id: 3, Name: "test3.jar", Status: "Archive File Within Another Archive"},
				{Id: 4, Name: "test4.jar", Status: "Archive File Within Another Archive"},
			},
			Issues: []report.Issue{},
		}

		nestedArchives(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Issues[0].Description, "2 nested archives")

		assert.Equal(t, 1, len(testReport.Recommendations))
	})
}
