package checks

import (
	"github.com/antfie/scan_health/v2/report"
	"testing"
)

func TestDependenciesSelected(t *testing.T) {
	t.Run("Unselected dependencies should not be reported", func(t *testing.T) {
		mockReport := &report.Report{
			Modules: []report.Module{
				{Name: "Test",
					Instances: []report.ModuleInstance{
						{IsDependency: true},
					}},
			},
		}

		dependenciesSelected(mockReport)

		if len(mockReport.Issues) != 0 {
			t.Errorf("Issues were reported which should not have been")
		}
	})

	t.Run("Selected dependencies should be reported", func(t *testing.T) {
		mockReport := &report.Report{
			Modules: []report.Module{
				{Name: "Test",
					Instances: []report.ModuleInstance{
						{IsDependency: true},
						{IsSelected: true},
					}},
			},
		}

		dependenciesSelected(mockReport)

		if len(mockReport.Issues) != 1 {
			t.Errorf("Issues were not reported which should have been")
		}
	})
}
