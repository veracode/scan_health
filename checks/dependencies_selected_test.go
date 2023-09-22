package checks

import (
	"github.com/antfie/scan_health/v2/report"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDependenciesSelected(t *testing.T) {
	t.Run("Unselected dependencies should not be reported", func(t *testing.T) {
		t.Parallel()
		mockReport := &report.Report{
			Modules: []report.Module{
				{Name: "Test",
					Instances: []report.ModuleInstance{
						{IsDependency: true},
					}},
			},
		}

		dependenciesSelected(mockReport)

		assert.Equal(t, len(mockReport.Issues), 0, "Issues were reported which should not have been")
	})

	t.Run("Selected dependencies should be reported", func(t *testing.T) {
		t.Parallel()
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

		assert.Equal(t, len(mockReport.Issues), 1, "Issues were not reported which should have been")
	})
}
