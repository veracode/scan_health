package checks

import "github.com/antfie/scan_health/v2/report"

func PerformChecks(r *report.Report) {
	// This is a priority ordered list of checks to perform so the customer sees the most important issues first

	// We put these at the top because other checks may rely on the output of these
	ignoreJunkFiles(r)
	thirdParty(r)

	// No flaws is a big deal
	flawCount(r)

	// Missing/unscannable components
	fatalErrors(r)
	unscannableJava(r)
	detectUnwantedFiles(r)
	nestedArchives(r)
	missingPrecompiledFiles(r)
	missingSCAComponents(r)
	unselectedJavaScriptModules(r)
	unexpectedSourceCode(r)
	missingSupportingFiles(r)
	missingDebugSymbols(r)
	unsupportedPlatformOrCompiler(r)
	gradleWrapper(r)

	// Undesirable things uploaded
	sensitiveFiles(r)
	repositories(r)
	nodeModules(r)
	testingArtefacts(r)
	tooManyFilesUploaded(r)
	excessMicrosoft(r)
	looseClassFiles(r)

	// Module selection issues
	unselectedFirstParty(r)
	dependenciesSelected(r)
	duplicateModules(r)
	moduleWarnings(r)

	// Others
	previousScan(r)
	minifiedJavaScript(r)
	releaseBuild(r)
	sizes(r)
	moduleCount(r)
	regularScans(r)

	// Finally
	generalRecommendations(r)
}
