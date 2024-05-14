package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veracode/scan_health/v2/report"
)

func TestGoWorkspaces(t *testing.T) {

	t.Run("Presence of go.work", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.exe", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "go.woRk", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		unsupportedGoWorkspaceFiles(&testReport)
		assert.NotEmpty(t, testReport.Issues)
	})

	t.Run("Presence of go.work.sum", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.exe", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "go.woRk.suM", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		unsupportedGoWorkspaceFiles(&testReport)
		assert.NotEmpty(t, testReport.Issues)
	})

	t.Run("No go workspace files", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.exe", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "go. woRk", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		unsupportedGoWorkspaceFiles(&testReport)
		assert.Empty(t, testReport.Issues)
	})
}
