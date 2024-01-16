package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veracode/scan_health/v2/report"
)

func TestUnselectedJavaScript(t *testing.T) {

	t.Run("No JavaScript modules", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "file3.exe",
					IsIgnored:    false,
					IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: true},
					}},
			},
		}

		unselectedJavaScriptModules(&mockReport)
		assert.Empty(t, mockReport.Issues)
	})

	t.Run("JS files within selected", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "js files within test.zip",
					IsIgnored:    false,
					IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: true},
					}},
			},
		}

		unselectedJavaScriptModules(&mockReport)
		assert.Empty(t, mockReport.Issues)
	})

	t.Run("Ignoring Node Modules", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "test_nodemodule_",
					IsIgnored:    false,
					IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: true},
					}},
			},
		}

		unselectedJavaScriptModules(&mockReport)
		assert.Empty(t, mockReport.Issues)
	})

	t.Run("Single JS files within not selected", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "js files within test.zip",
					IsIgnored:    false,
					IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: false},
					}},
			},
		}

		unselectedJavaScriptModules(&mockReport)
		if !assert.Equal(t, 1, len(mockReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, mockReport.Issues[0].Description, "A JavaScript module was not selected")
		assert.Equal(t, report.IssueSeverityMedium, mockReport.Issues[0].Severity)
		assert.Equal(t, 2, len(mockReport.Recommendations))
	})

	t.Run("Multiple JS files within not selected", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "js files within test.zip",
					IsIgnored:    false,
					IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: false},
					}},
				{Name: "js files extracted from test2.war",
					IsIgnored:    false,
					IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: false},
					}},
			},
		}

		unselectedJavaScriptModules(&mockReport)
		if !assert.Equal(t, 1, len(mockReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, mockReport.Issues[0].Description, "2 JavaScript module")
		assert.Equal(t, report.IssueSeverityMedium, mockReport.Issues[0].Severity)
		assert.Equal(t, 2, len(mockReport.Recommendations))
	})

	t.Run("Ignore Fatal JS files", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				{Name: "js files within test.zip",
					IsIgnored:    false,
					IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: false},
						{HasFatalErrors: true},
					}},
				{Name: "js files extracted from test2.war",
					IsIgnored:    false,
					IsThirdParty: false,
					Instances: []report.ModuleInstance{
						{IsDependency: false},
						{IsSelected: false},
						{HasFatalErrors: true},
					}},
			},
		}

		unselectedJavaScriptModules(&mockReport)
		assert.Empty(t, mockReport.Issues)
	})

}
