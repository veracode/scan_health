package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veracode/scan_health/v2/report"
)

func TestUnscannableJava(t *testing.T) {

	t.Run("Normal Java File", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.jar",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
					}},
			},
		}

		unscannableJava(&mockReport)

		assert.Empty(t, mockReport.Issues)
	})

	t.Run("No Java Files", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.exe",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
					}},
			},
		}

		unscannableJava(&mockReport)

		assert.Empty(t, mockReport.Issues)
	})

	t.Run("2 Java Files, 1 unscannable", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.exe",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
					}},
				{Name: "file4.jar",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{HasFatalErrors: true},
						{Issues: []string{`A Java module "file4.jar" was not scannable as it contained no Java binaries.`}},
					},
				},
			},
		}

		unscannableJava(&mockReport)

		assert.Equal(t, 1, len(mockReport.Issues))
		assert.Equal(t, 1, len(mockReport.Recommendations))
	})

	t.Run("3 Java Files, 2 unscannable", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.exe",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
					}},
				{Name: "file4.jar",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{HasFatalErrors: true},
						{Issues: []string{`A Java module "file4.jar" was not scannable as it contained no Java binaries.`}},
					},
				},
				{Name: "file5.war",
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{HasFatalErrors: true},
						{Issues: []string{`A Java module "file5.war" was not scannable as it contained no Java binaries.`}},
					},
				},
			},
		}

		unscannableJava(&mockReport)

		if !assert.Equal(t, 1, len(mockReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, mockReport.Issues[0].Description, "2 Java modules")
		assert.Equal(t, 1, len(mockReport.Recommendations))
	})
}
