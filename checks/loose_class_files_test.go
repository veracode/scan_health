package checks

import (
	"testing"

	"github.com/antfie/scan_health/v2/report"
	"github.com/stretchr/testify/assert"
)

// Test Cases
func TestLooseClassFiles(t *testing.T) {

	// Test Case 1: No Roslyn Files
	t.Run("Normal Files", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.exe", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		looseClassFiles(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	t.Run("One Class File", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.class", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		looseClassFiles(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Issues[0].Description, "Java class files were not packaged")
		assert.Equal(t, 1, len(testReport.Recommendations))
	})

	t.Run("One class Module", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "Class Files Within test.zip",
					IsThirdParty: false,
					IsIgnored:    false,
					Instances: []report.ModuleInstance{{
						IsSelected:     true,
						HasFatalErrors: false,
					}},
				},
			},

			Issues: []report.Issue{},
		}

		looseClassFiles(&testReport)
		assert.Contains(t, testReport.Issues[0].Description, "Java class files were not packaged")
		assert.Equal(t, 1, len(testReport.Recommendations))
	})

	t.Run("A Class Files and a class module", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "Class Files Within test.zip",
					IsThirdParty: false,
					IsIgnored:    false,
					Instances: []report.ModuleInstance{{
						IsSelected:     true,
						HasFatalErrors: false,
					}},
				},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.class", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 333333, Name: "test2.jar", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},

			Issues: []report.Issue{},
		}

		looseClassFiles(&testReport)
		assert.Contains(t, testReport.Issues[0].Description, "Java class files were not packaged")
		assert.Equal(t, 1, len(testReport.Recommendations))
	})
}
