package comparisons

import (
	"fmt"
	"github.com/veracode/scan_health/v2/report"
	"github.com/veracode/scan_health/v2/utils"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func reportFlawDifferences(a, b *report.Report) {
	flawsFromA := getFlawsFromReport(a)
	flawsFromB := getFlawsFromReport(b)

	reportFlawStateDifferences(flawsFromA, flawsFromB)
	reportFlawMitigationDifferences(flawsFromA, flawsFromB)
	reportFlawLineNumberChanges(flawsFromA, flawsFromB)
	reportPolicyAffectingFlawDifferences(flawsFromA, flawsFromB)
	reportNonPolicyAffectingFlawDifferences(flawsFromA, flawsFromB)
	reportClosedFlawDifferences(flawsFromA, flawsFromB)
}

func reportFlawStateDifferences(flawsFromA, flawsFromB []report.FlawDetails) {
	var r strings.Builder

	compareFlawStates(&r, flawsFromA, flawsFromB)

	if r.Len() > 0 {
		utils.PrintTitle("Flaw State Differences")
		utils.ColorPrintf(r.String())
	}
}

func reportFlawMitigationDifferences(flawsFromA, flawsFromB []report.FlawDetails) {
	var r strings.Builder

	compareFlawMitigations(&r, flawsFromA, flawsFromB)

	if r.Len() > 0 {
		utils.PrintTitle("Flaw Mitigation Differences")
		utils.ColorPrintf(r.String())
	}
}

func reportFlawLineNumberChanges(flawsFromA, flawsFromB []report.FlawDetails) {
	var r strings.Builder

	compareFlawLineNumberChanges(&r, flawsFromA, flawsFromB)

	if r.Len() > 0 {
		utils.PrintTitle("Flaw Line Number Differences")
		utils.ColorPrintf(r.String())
	}
}

func reportPolicyAffectingFlawDifferences(flawsFromA, flawsFromB []report.FlawDetails) {
	var r strings.Builder

	compareFlaws(&r, "A", flawsFromA, flawsFromB, true, false)
	compareFlaws(&r, "B", flawsFromA, flawsFromB, true, false)

	if r.Len() > 0 {
		utils.PrintTitle("Policy Affecting Open Flaw Differences")
		utils.ColorPrintf(r.String())
	}
}

func reportNonPolicyAffectingFlawDifferences(flawsFromA, flawsFromB []report.FlawDetails) {
	var r strings.Builder

	compareFlaws(&r, "A", flawsFromA, flawsFromB, false, false)
	compareFlaws(&r, "B", flawsFromA, flawsFromB, false, false)

	if r.Len() > 0 {
		utils.PrintTitle("Non Policy Affecting Open Flaw Differences")
		utils.ColorPrintf(r.String())
	}
}

func reportClosedFlawDifferences(flawsFromA, flawsFromB []report.FlawDetails) {
	var r strings.Builder

	compareFlaws(&r, "A", flawsFromA, flawsFromB, false, true)
	compareFlaws(&r, "B", flawsFromA, flawsFromB, false, true)

	if r.Len() > 0 {
		utils.PrintTitle("Closed Flaw Differences")
		utils.ColorPrintf(r.String())
	}
}

func getFlawsFromReport(r *report.Report) []report.FlawDetails {
	var flaws []report.FlawDetails

	for _, module := range r.Modules {
		flaws = append(flaws, module.FlawDetails...)
	}

	return flaws
}

func getSortedCwes(flaws []report.FlawDetails) []int {
	var cwes []int
	for _, thisSideFlaw := range flaws {
		if !utils.IsInIntArray(thisSideFlaw.CWE, cwes) {
			cwes = append(cwes, thisSideFlaw.CWE)
		}
	}

	sort.Ints(cwes[:])
	return cwes
}

func isFlawNumberInArray(flawId int, flaws []report.FlawDetails) bool {
	for _, flaw := range flaws {
		if flawId == flaw.ID {
			return true
		}
	}

	return false
}

func compareFlaws(r *strings.Builder, side string, flawsFromA, flawsFromB []report.FlawDetails, policyAffecting bool, onlyClosed bool) {
	flawsInThisReport := flawsFromA
	flawsInOtherReport := flawsFromB

	if side == "B" {
		flawsInThisReport = flawsFromB
		flawsInOtherReport = flawsFromA
	}

	for _, cwe := range getSortedCwes(flawsInThisReport) {
		var flawsOnlyInThisScan []int

		for _, thisSideFlaw := range flawsInThisReport {
			if thisSideFlaw.CWE != cwe {
				continue
			}

			if onlyClosed {
				if thisSideFlaw.IsOpen() {
					continue
				}
			} else {
				if policyAffecting && !(thisSideFlaw.IsOpen() && thisSideFlaw.AffectsPolicyCompliance) {
					continue
				}

				if !policyAffecting && !(thisSideFlaw.IsOpen() && !thisSideFlaw.AffectsPolicyCompliance) {
					continue
				}
			}

			if !isFlawNumberInArray(thisSideFlaw.ID, flawsInOtherReport) {
				flawsOnlyInThisScan = append(flawsOnlyInThisScan, thisSideFlaw.ID)
			}
		}

		if len(flawsOnlyInThisScan) > 0 {
			r.WriteString(fmt.Sprintf("%s: %dx CWE-%d = %s\n",
				utils.GetFormattedOnlyInSideString(side),
				len(flawsOnlyInThisScan),
				cwe,
				getSortedIntArrayAsFormattedString(flawsOnlyInThisScan)))
		}
	}
}

func getSortedIntArrayAsFormattedString(list []int) string {
	sort.Ints(list[:])

	var output []string
	for _, x := range list {
		output = append(output, strconv.Itoa(x))
	}

	return strings.Join(output, ",")
}

func compareFlawStates(r *strings.Builder, flawsFromA, flawsFromB []report.FlawDetails) {
	stateChanges := make(map[string][]int)

	for _, aFlaw := range flawsFromA {
		for _, bFlaw := range flawsFromB {
			if aFlaw.ID != bFlaw.ID {
				continue
			}

			if aFlaw.RemediationStatus == bFlaw.RemediationStatus {
				continue
			}

			var stateChange = fmt.Sprintf("%s %-9s => %s %-9s: CWE-%d",
				utils.GetFormattedSideString("A"),
				aFlaw.RemediationStatus,
				utils.GetFormattedSideString("B"),
				bFlaw.RemediationStatus,
				aFlaw.CWE)

			stateChanges[stateChange] = append(stateChanges[stateChange], aFlaw.ID)
		}
	}

	sortedKeys := make([]string, 0, len(stateChanges))
	for k := range stateChanges {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		var flawIds = stateChanges[key]

		var formattedS = strings.Replace(key, "CWE", fmt.Sprintf("%dx CWE", len(flawIds)), 1)
		r.WriteString(fmt.Sprintf("%s = %s\n", formattedS, getSortedIntArrayAsFormattedString(flawIds)))
	}
}

func compareFlawMitigations(r *strings.Builder, flawsFromA, flawsFromB []report.FlawDetails) {
	for _, aFlaw := range flawsFromA {
		for _, bFlaw := range flawsFromB {
			if aFlaw.ID != bFlaw.ID {
				continue
			}

			if aFlaw.MitigationStatus != bFlaw.MitigationStatus {
				r.WriteString(fmt.Sprintf("%d (CWE-%d): %s: %s, %s: %s\n",
					aFlaw.ID,
					aFlaw.CWE,
					utils.GetFormattedSideString("A"),
					cases.Title(language.English).String(aFlaw.MitigationStatus),
					utils.GetFormattedSideString("B"),
					cases.Title(language.English).String(bFlaw.MitigationStatus)))
			}
		}
	}
}

func compareFlawLineNumberChanges(r *strings.Builder, flawsFromA, flawsFromB []report.FlawDetails) {
	for _, aFlaw := range flawsFromA {
		for _, bFlaw := range flawsFromB {
			if aFlaw.ID != bFlaw.ID {
				continue
			}

			if aFlaw.LineNumber != bFlaw.LineNumber {
				r.WriteString(fmt.Sprintf("%d (CWE-%d): %s: %d, %s: %d\n",
					aFlaw.ID,
					aFlaw.CWE,
					utils.GetFormattedSideString("A"),
					aFlaw.LineNumber,
					utils.GetFormattedSideString("B"),
					bFlaw.LineNumber))
			}
		}
	}
}
