package checks

import (
	"testing"

	"github.com/antfie/scan_health/v2/report"
	"github.com/stretchr/testify/assert"
)

func TestMinifiedJavaScript(t *testing.T) {

	t.Run("No Issues", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "Test",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
					}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.js", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2.js", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
		}

		minifiedJavaScript(&mockReport)

		assert.Equal(t, len(mockReport.Issues), 0)
	})

	t.Run("No JS Minified Issues", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "Test",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{Issues: []string{"No PDB files found"}},
					}},
			},
		}

		minifiedJavaScript(&mockReport)

		assert.Equal(t, len(mockReport.Issues), 0)
	})

	t.Run("Testing for JS Minified Warning", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "JS files within Test",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{Issues: []string{"because we think it is minified"}},
					}},
				{Name: "JS files within Test2",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{Issues: []string{"because we think it is minified"}},
					}},
			},
		}

		minifiedJavaScript(&mockReport)

		if !assert.Equal(t, len(mockReport.Issues), 1) {
			t.FailNow()
		}

		assert.Contains(t, mockReport.Issues[0].Description, "2 minified")
		assert.Equal(t, mockReport.Issues[0].Severity, report.IssueSeverityMedium)
		assert.Equal(t, len(mockReport.Recommendations), 2)
	})

	t.Run("Testing for /dist/ JS Files", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "JS files within Test",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{Issues: []string{"Random issue in /dist/blah.js"}},
					}},
			},
		}

		minifiedJavaScript(&mockReport)

		assert.Equal(t, len(mockReport.Issues), 1)
		assert.Equal(t, mockReport.Issues[0].Severity, report.IssueSeverityMedium)
		assert.Equal(t, len(mockReport.Recommendations), 2)
	})

	t.Run("Testing for minification by name", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.min.js", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2.js", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file3.min.js", MD5: "hash3", IsIgnored: false, IsThirdParty: false},
			},
		}

		minifiedJavaScript(&mockReport)
		if !assert.Equal(t, len(mockReport.Issues), 1) {
			t.FailNow()
		}

		assert.Contains(t, mockReport.Issues[0].Description, "2 minified")
		assert.Equal(t, mockReport.Issues[0].Severity, report.IssueSeverityMedium)

		assert.Equal(t, len(mockReport.Recommendations), 2)
	})

}
