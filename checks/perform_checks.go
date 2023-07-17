package checks

import "github.com/antfie/scan_health/v2/report"

func PerformChecks(r *report.Report) {
	// This is a priority ordered list of checks to perform so the customer sees the most important issues first

	// We put these at the top because other checks may rely on the output of these
	ignoreJunkFiles(r)
	thirdParty(r)

	// Missing/unscannable components
	fatalErrors(r)
	detectUnwantedFiles(r)
	nestedArchives(r)
	unselectedFirstParty(r)
	unselectedJavaScriptModules(r)
	missingPrecompiledFiles(r)
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
	dependenciesSelected(r)
	duplicateModules(r)
	moduleWarnings(r)

	// Others
	previousScan(r)
	minifiedJavaScript(r)
	releaseBuild(r)
	analysisSize(r)
	moduleCount(r)
	regularScans(r)
	flawCount(r)

	// Finally
	generalRecommendations(r)
}

// TODO - Backlog of checks to add
// * There were apparent external-facing application components that had not been selected as entry points for analysis. This could result in the reduced scan coverage.
// * Consider an application profile for each supported version of the application in production so the security team can see the risk of each specific version.
// * Be sure to include all the components that make up the application within the upload. Do not omit any second or third-party libraries from the upload.
// * It was observed that there were several sandboxes with names that suggest the team uses a sandbox for each significant version of the application. Further there were sandboxes with names that suggest different components of the application were being scanned in each e.g. "TODO", "TODO". The security team will expect the policy scan to contain all the components of the application to get a complete picture of all the risk. Since we can only promote one sandbox at a time to the policy level there is a concern that what is promoted to the policy level is not the entire application.
// * There were apparent external-facing application components (“TODO”, “TODO”) that had not been selected as entry points for analysis. This could result in the reduced scan coverage.
// * Support Issue: The image \"X\" contains statically linked standard libraries. Proceeding with these libraries included will degrade the performance of the analysis and quality of the results. Disable static linking by omitting -static and -static-libstdc++ GCC options.",
// * Java or .NET Files Compiled without Debug Symbols - 1 File
// * obfuscation - Compiled using obfuscation &#x28;Optional&#x29;

// Bugs:
// * https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:17521:293966:26537614:26508585:26524235::::341599 - duplicate module, one as dependency, file sizes differ, both selected
// * dependencies were selected here https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:17521:293966:26537614:26508585:26524235::::341599
// * Fix "and 1 others"
