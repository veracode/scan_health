package checks

import (
	"testing"

	"github.com/antfie/scan_health/v2/report"
	"github.com/stretchr/testify/assert"
)

// Test Cases
func TestMissingPrecompiledFiles(t *testing.T) {

	t.Run("Valid Files and Modules", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "file1.dll", Instances: []report.ModuleInstance{{IsDependency: false}}},
				{Name: "file2.dll", Instances: []report.ModuleInstance{{IsDependency: false}}},
				{Name: "file3.dll", Instances: []report.ModuleInstance{{IsDependency: false}}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.dll", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file3.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		missingPrecompiledFiles(&testReport)
		assert.Empty(t, testReport.Issues)
	})

	t.Run("CSHTML File uploaded", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "file1.dll", Instances: []report.ModuleInstance{{IsDependency: false}}},
				{Name: "file2.dll", Instances: []report.ModuleInstance{{IsDependency: false}}},
				{Name: "file3.dll", Instances: []report.ModuleInstance{{IsDependency: false}}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.dll", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 333333, Name: "login.cshtml", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 444444, Name: "file3.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		missingPrecompiledFiles(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Issues[0].Description, "A .NET view/template/control file")
		assert.Equal(t, 1, len(testReport.Recommendations))
	})

	t.Run("CSHTML and ASPX Files uploaded", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "file1.dll", Instances: []report.ModuleInstance{{IsDependency: false}}},
				{Name: "file2.dll", Instances: []report.ModuleInstance{{IsDependency: false}}},
				{Name: "file3.dll", Instances: []report.ModuleInstance{{IsDependency: false}}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.dll", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 333333, Name: "login.cshtml", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 444444, Name: "logout.aspx", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 555555, Name: "file3.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		missingPrecompiledFiles(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Issues[0].Description, "2 .NET views/templates/control files")
		assert.Equal(t, 1, len(testReport.Recommendations))
	})

	t.Run("Module Reports No Precompiled Files", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "file1.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: true}}},
				{Name: "file2.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: true}}},
				{Name: "file3.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: true},
					{Issues: []string{"No precompiled files were found for this .NET web application"}}},
				},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.dll", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file3.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		missingPrecompiledFiles(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Issues[0].Description, "A .NET component was identified ")
	})

	t.Run("Multiple Modules Report No Precompiled Files", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "file1.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: true}}},
				{Name: "file2.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: true},
					{Issues: []string{"No precompiled files were found for this .NET web application"}}},
				},
				{Name: "file3.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: true},
					{Issues: []string{"No precompiled files were found for this .NET web application"}}},
				},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.dll", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file3.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		missingPrecompiledFiles(&testReport)
		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Issues[0].Description, "2 .NET components were identified ")
	})

	t.Run("Unselected Module With No Precompiled Files", func(t *testing.T) {
		t.Parallel()
		testReport := report.Report{
			Modules: []report.Module{
				{Name: "file1.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: true}}},
				{Name: "file2.dll", Instances: []report.ModuleInstance{{IsDependency: false}, {IsSelected: false},
					{Issues: []string{"No precompiled files were found for this .NET web application"}}},
				},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "file1.dll", MD5: "hash1", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file2.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
				{Id: 222222, Name: "file3.dll", MD5: "hash2", IsIgnored: false, IsThirdParty: false},
			},
			Issues: []report.Issue{},
		}

		missingPrecompiledFiles(&testReport)
		assert.Empty(t, len(testReport.Issues))
	})
}
