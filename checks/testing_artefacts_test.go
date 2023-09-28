package checks

import (
	"testing"

	"github.com/antfie/scan_health/v2/report"
	"github.com/stretchr/testify/assert"
)

func TestTestingArtefacts(t *testing.T) {

	t.Run("No Testing Artefacts", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.jar",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
					}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 333333, Name: "file3.jar", MD5: "hash2", IsIgnored: false},
			},
		}

		testingArtefacts(&mockReport)

		assert.Empty(t, mockReport.Issues)
	})

	t.Run("Testing File uploaded", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.dll",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
					}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "Antlr3.Runtime.dll", MD5: "hash1", IsIgnored: false},
				{Id: 222222, Name: "file3.unittests.dll", MD5: "hash2", IsIgnored: false},
				{Id: 222222, Name: "file15.unittest.dll", MD5: "hash2", IsIgnored: false},
			},
		}

		testingArtefacts(&mockReport)

		if !assert.Equal(t, len(mockReport.Issues), 1) {
			t.FailNow()
		}

		assert.Contains(t, mockReport.Issues[0].Description, "2 testing artefacts")
		assert.Equal(t, 1, len(mockReport.Recommendations))

	})

	t.Run("Testing Module found, not selected", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "moq.dll",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
					}},
				{Name: "standalone.unittests.dll",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
					}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "upload.zip", MD5: "hash1", IsIgnored: false},
			},
		}

		testingArtefacts(&mockReport)

		assert.Empty(t, mockReport.Issues)
	})

	t.Run("Testing Module found, selected", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "moq.dll",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: true},
					}},
				{Name: "standalone.unittests.dll",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: true},
					}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "upload.zip", MD5: "hash1", IsIgnored: false},
			},
		}

		testingArtefacts(&mockReport)

		if !assert.Equal(t, 1, len(mockReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, mockReport.Issues[0].Description, "2 testing artefacts")
		assert.Equal(t, 2, len(mockReport.Recommendations))
	})

	t.Run("Module dependant on testing artefacts", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "MyApp.exe",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: true},
						{Issues: []string{"something including test/dependency.dll"}},
					}},
				{Name: "AnotherApp.dll",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: true},
						{Issues: []string{"something including test/userDependency.dll"}},
					}},
			},
			UploadedFiles: []report.UploadedFile{
				{Id: 111111, Name: "upload.zip", MD5: "hash1", IsIgnored: false},
			},
		}

		testingArtefacts(&mockReport)

		if !assert.Equal(t, 1, len(mockReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, mockReport.Issues[0].Description, "2 modules")
		assert.Equal(t, 1, len(mockReport.Recommendations))
	})
}
