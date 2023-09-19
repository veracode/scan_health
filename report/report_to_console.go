package report

import (
	"fmt"
	"github.com/antfie/scan_health/v2/utils"
	"github.com/fatih/color"
	"strings"
	"time"
)

func (r *Report) renderToConsole() {
	renderScanSummaryToConsole(r)
	renderUploadedFilesToConsole(r)
	renderSelectedModulesToConsole(r)
	renderFlawSummaryToConsole(r.Flaws)
	renderIssues(r.Issues)
	renderRecommendations(r.Recommendations)
}

func renderScanSummaryToConsole(report *Report) {
	utils.PrintTitle("Scan Summary")

	if report.Scan.BusinessUnit != "Not Specified" {
		fmt.Printf("Business unit:      %s\n", report.Scan.BusinessUnit)
	}

	fmt.Printf("Application:        %s\n", report.Scan.ApplicationName)

	if len(report.Scan.SandboxName) > 0 {
		fmt.Printf("Sandbox:            %s\n", report.Scan.SandboxName)
	}

	fmt.Printf("Scan name:          %s\n", report.Scan.ScanName)
	fmt.Printf("Review modules URL: %s\n", report.Scan.ReviewModulesUrl)
	fmt.Printf("Triage flaws URL:   %s\n", report.Scan.TriageFlawsUrl)
	fmt.Printf("Files uploaded:     %d\n", len(report.UploadedFiles))
	fmt.Printf("Total modules:      %d\n", len(report.Modules))
	fmt.Printf("Modules selected:   %d\n", len(report.GetSelectedModules()))
	fmt.Printf("Engine version:     %s (Release notes: https://docs.veracode.com/updates/r/c_all_static)\n", report.Scan.EngineVersion)
	fmt.Printf("Submitted:          %s (%s ago)\n", report.Scan.SubmittedDate, utils.FormatDuration(time.Since(report.Scan.SubmittedDate)))
	fmt.Printf("Published:          %s (%s ago)\n", report.Scan.PublishedDate, utils.FormatDuration(time.Since(report.Scan.PublishedDate)))
	fmt.Printf("Duration:           %s\n", utils.FormatDuration(report.Scan.ScanDuration))

	if !report.Scan.IsLatestScan && report.LastAppActivity.Year() > 2000 {
		fmt.Printf("Latest app scan:    %s (%s ago)\n", report.LastAppActivity, utils.FormatDuration(time.Since(report.LastAppActivity)))
	}
}

func renderFlawSummaryToConsole(flaws FlawSummary) {
	utils.PrintTitle("Flaw Summary")

	total := fmt.Sprintf("Total (less fixed):         %d\n", flaws.Total-flaws.Fixed)

	if flaws.Total == 0 || flaws.Total > utils.MaximumFlawCountThreshold {
		color.HiYellow(total)
	} else {
		fmt.Print(total)
	}

	fmt.Printf("Fixed (no longer reported): %d\n", flaws.Fixed)
	fmt.Printf("Policy affecting:           %d\n", flaws.TotalAffectingPolicy)
	fmt.Printf("Mitigated:                  %d\n", flaws.Mitigated)
	fmt.Printf("Open affecting policy:      %d\n", flaws.OpenAffectingPolicy)
	fmt.Printf("Open not affecting policy:  %d\n", flaws.OpenButNotAffectingPolicy)
}

func renderUploadedFilesToConsole(report *Report) {
	// No point in listing every file if there are many
	if len(report.UploadedFiles) > 25 {
		return
	}

	utils.PrintTitle("Uploaded files")

	for _, uploadedFile := range report.UploadedFiles {
		fmt.Printf("* %s (MD5: %s)\n", uploadedFile.Name, uploadedFile.MD5)
	}
}

func renderSelectedModulesToConsole(report *Report) {
	// No point in listing every selected module if there are many
	if len(report.GetSelectedModules()) > 25 {
		return
	}

	utils.PrintTitle("Modules Selected for Analysis")

	for _, selectedModule := range report.GetSelectedModules() {
		var notes []string

		for _, moduleInstance := range selectedModule.Instances {
			if len(moduleInstance.MD5) > 0 {
				notes = append(notes, fmt.Sprintf("MD5: %s", moduleInstance.MD5))
			}

			if len(moduleInstance.OperatingSystem) > 0 {
				notes = append(notes, fmt.Sprintf("OS: %s", moduleInstance.OperatingSystem))
			}

			if len(moduleInstance.Architecture) > 0 {
				notes = append(notes, fmt.Sprintf("Architecture: %s", moduleInstance.Architecture))
			}

			if len(moduleInstance.Compiler) > 0 {
				notes = append(notes, fmt.Sprintf("Compiler: %s", moduleInstance.Compiler))
			}

			if len(moduleInstance.Platform) > 0 {
				notes = append(notes, fmt.Sprintf("Platform: %s", moduleInstance.Platform))
			}
		}

		additionalNotes := ""

		if len(notes) > 0 {
			additionalNotes = " (" + strings.Join(notes, ", ") + ")"
		}

		fmt.Printf("* %s%s\n", selectedModule.Name, additionalNotes)
	}
}

func (issue Issue) getIcon() string {
	if issue.Severity == IssueSeverityMedium {
		return iconWarning
	}

	return iconError
}

func renderIssues(issues []Issue) {
	if len(issues) < 1 {
		return
	}

	utils.PrintTitle("Issues")

	// Render the high severity issues first
	for _, issue := range issues {
		if issue.Severity == IssueSeverityHigh {
			utils.ColorPrintf(fmt.Sprintf("%s %s\n", issue.getIcon(), issue.Description))
		}
	}

	for _, issue := range issues {
		if issue.Severity != IssueSeverityHigh {
			utils.ColorPrintf(fmt.Sprintf("%s %s\n", issue.getIcon(), issue.Description))
		}
	}
}

func renderRecommendations(recommendations []string) {
	if len(recommendations) < 1 {
		return
	}

	utils.PrintTitle("Recommendations")

	for _, recommendation := range recommendations {
		utils.ColorPrintf(fmt.Sprintf("%s %s\n", iconRecommendation, recommendation))
	}
}
