package checks

import (
	"strings"
	"testing"

	"github.com/antfie/scan_health/v2/report"
	"github.com/stretchr/testify/assert"
)

// Test Cases
func TestFilesToIgnore(t *testing.T) {

	// Test Case 1: No duplicates
	t.Run("No Files to Ignore", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "this is a valid file.jar", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "this is a valid file.exe", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		ignoreJunkFiles(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	t.Run("Ignore PDB and .gitignore as Special Cases", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: ".gitignore", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "test.pdb", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		ignoreJunkFiles(&testReport)
		assert.Empty(t, testReport.Issues)
		assert.True(t, testReport.UploadedFiles[1].IsIgnored)
	})

	t.Run("1 file to Ignore", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "web.config", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "this is a valid file.exe", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		ignoreJunkFiles(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.True(t, strings.Contains(testReport.Issues[0].Description, "unnecessary"))

		if !assert.Equal(t, 1, len(testReport.Recommendations)) {
			t.FailNow()
		}
	})

	t.Run("2 files to Ignore", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "web.config", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "Makefile", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 333333, Name: "Test.exe", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		ignoreJunkFiles(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.True(t, strings.Contains(testReport.Issues[0].Description, "2 unnecessary"))

		if !assert.Equal(t, 1, len(testReport.Recommendations)) {
			t.FailNow()
		}
	})

	t.Run("Any modules derived from ignored files should also be ignored", func(t *testing.T) {
		t.Parallel()

		testReport := report.Report{
			Modules: []report.Module{
				{Name: "dist.js.map", Instances: []report.ModuleInstance{}},
			},
		}

		ignoreJunkFiles(&testReport)

		assert.True(t, testReport.Modules[0].IsIgnored)
	})
}
