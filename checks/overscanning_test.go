package checks

import (
	"strings"
	"testing"

	"github.com/antfie/scan_health/v2/report"
	"github.com/stretchr/testify/assert"
)

func TestOverScanning(t *testing.T) {
	t.Run("Over-scanning detected", func(t *testing.T) {
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
					Name: "app.dll", Instances: []report.ModuleInstance{
						{IsSelected: true},
					},
				},
			},
		}

		overScanning(&testReport)

		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.True(t, strings.Contains(testReport.Issues[0].Description, "This is because it was already included in the analysis"))

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
						{IsSelected: true},
					},
					DependencyOf: []string{"app.dll"},
				},
				{
					Name: "common-models.dll", Instances: []report.ModuleInstance{
						{IsSelected: true},
					},
					DependencyOf: []string{"app.dll"},
				},
				{
					Name: "app.dll", Instances: []report.ModuleInstance{
						{IsSelected: true},
					},
				},
			},
		}

		overScanning(&testReport)

		if !assert.Equal(t, 1, len(testReport.Issues)) {
			t.FailNow()
		}

		assert.True(t, strings.Contains(testReport.Issues[0].Description, "This is because they were already included in the analysis"))

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

		if !assert.Equal(t, 0, len(testReport.Issues)) {
			t.FailNow()
		}

		if !assert.Equal(t, 0, len(testReport.Recommendations)) {
			t.FailNow()
		}
	})

}
