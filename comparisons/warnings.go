package comparisons

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"strings"
	"time"
)

func reportWarnings(a, b *report.Report) {
	var r strings.Builder

	if a.Scan.EngineVersion != b.Scan.EngineVersion {
		r.WriteString("* The scan engine versions are different. This means there has been one or more deployments of the Veracode scan engine between these scans. This can sometimes explain why new flaws might be reported (due to improved scan coverage), and others are no longer reported (due to a reduction of False Positives)\n")
	}

	if a.Scan.AccountId != b.Scan.AccountId {
		r.WriteString("* These scans are from different accounts\n")
	}

	if a.Scan.ApplicationId != b.Scan.ApplicationId {
		r.WriteString("* These scans are from different application profiles\n")
	}

	if !a.Scan.IsLatestScan {
		if len(a.Scan.SandboxName) > 0 {
			r.WriteString(fmt.Sprintf("* Scan A is not the latest SAST scan in the sandbox \"%s\"\n", a.Scan.SandboxName))
		} else {
			r.WriteString("* Scan A is not the latest SAST policy scan\n")
		}
	}

	if !b.Scan.IsLatestScan {
		if len(b.Scan.SandboxName) > 0 {
			r.WriteString(fmt.Sprintf("* Scan B is not the latest SAST scan in the sandbox \"%s\"\n", b.Scan.SandboxName))
		} else {
			r.WriteString("* Scan B is not the latest SAST policy scan\n")
		}
	}

	const hoursInThirtyDays = 30 * 24

	if time.Since(a.Scan.SubmittedDate).Hours() >= hoursInThirtyDays && time.Since(b.Scan.SubmittedDate).Hours() >= hoursInThirtyDays {
		r.WriteString("* Both scans are older than 30 days. This means the files will have been deleted and Veracode support therefore require a newer scan to investigate any issues further.\n")
	} else if time.Since(a.Scan.SubmittedDate).Hours() >= hoursInThirtyDays {
		r.WriteString("* Scan A is older than 30 days. This means the files will have been deleted and Veracode support therefore require a newer scan to investigate any issues further.\n")
	} else if time.Since(b.Scan.SubmittedDate).Hours() >= hoursInThirtyDays {
		r.WriteString("* Scan B is older than 30 days. This means the files will have been deleted and Veracode support therefore require a newer scan to investigate any issues further.\n")
	}

	if r.Len() > 0 {
		utils.PrintTitle("Warnings")
		color.HiYellow(r.String())
	}
}
