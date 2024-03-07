package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veracode/scan_health/v2/report"
)

func TestOverScanning(t *testing.T) {
	t.Run("Over-scanning detected", func(t *testing.T) {
		t.Parallel()

		testReport := report.Report{
			Modules: []report.Module{
				{
					Name: "common-lib.dll", Instances: []report.ModuleInstance{
						{IsSelected: true, Source: report.DetailedReportModuleSelected},
					},
					DependencyOf: []string{"app.dll"},
				},
				{
					Name: "app.dll", Instances: []report.ModuleInstance{
						{IsSelected: true, Source: report.DetailedReportModuleSelected},
					},
				},
			},
		}

		overScanning(&testReport)

		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Issues[0].Description, "This is because it was already included in the analysis")

		if !assert.Equal(t, 2, len(testReport.Recommendations)) {
			t.FailNow()
		}
	})

	t.Run("Multiple instances of over scanning detected", func(t *testing.T) {
		t.Parallel()

		testReport := report.Report{
			Modules: []report.Module{
				{
					Name: "common-lib.dll", Instances: []report.ModuleInstance{
						{IsSelected: true, Source: report.DetailedReportModuleSelected},
					},
					DependencyOf: []string{"app.dll"},
				},
				{
					Name: "common-models.dll", Instances: []report.ModuleInstance{
						{IsSelected: true, Source: report.DetailedReportModuleSelected},
					},
					DependencyOf: []string{"app.dll"},
				},
				{
					Name: "app.dll", Instances: []report.ModuleInstance{
						{IsSelected: true, Source: report.DetailedReportModuleSelected},
					},
				},
			},
		}

		overScanning(&testReport)

		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, testReport.Issues[0].Description, "This is because they were already included in the analysis")

		if !assert.Equal(t, 2, len(testReport.Recommendations)) {
			t.FailNow()
		}
	})

	t.Run("There is no over-scanning because the consuming library was not selected for analysis", func(t *testing.T) {
		t.Parallel()

		testReport := report.Report{
			Modules: []report.Module{
				{
					Name: "common-lib.dll", Instances: []report.ModuleInstance{
						{IsSelected: true},
					},
					DependencyOf: []string{"app.dll"},
				},
				{
					Name: "app.dll", Instances: []report.ModuleInstance{},
				},
			},
		}

		overScanning(&testReport)

		if !assert.Empty(t, testReport.Issues) {
			t.FailNow()
		}

		assert.Empty(t, testReport.Recommendations)
	})

}
