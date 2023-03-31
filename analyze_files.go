package main

import (
	"sort"
	"strings"
)

func (data Data) analyzeUploadedFiles() {
	var report strings.Builder

	var files []string

	for _, file := range data.PrescanFileList.Files {
		files = append(files, file.Name)

	}

	sort.Strings(files[:])

	if len(files) > 10000 {
		report.WriteString(formatWarningStringFormat(
			"%d files were present. This is a lot of files which is usually an indicator that something is not correct\n",
			len(files)))
	}

	detectSensitiveFiles(data, &report, files)
	detectTestArtefacts(data, &report, files)
	detectNodeModules(data, &report)
	detectRoslyn(data, &report, files)
	detectGit(data, &report, files)
	detectUnwantedFiles(data, &report, files, ".zip", "nested zip file", []string{"Do not upload archives (nested archives) within the upload package"})
	detectUnwantedFiles(data, &report, files, ".7z", "7-zip file", []string{"Veracode does not support 7-zip. Consider zip files instead"})
	detectUnwantedFiles(data, &report, files, ".java", "Java source code file", []string{"Do not upload Java source code files. They will not be scanned", "Veracode requires Java application to be compiled into a .jar, .war or .ear file"})
	detectUnwantedFiles(data, &report, files, ".class", "Java class file", []string{"Do not upload Java class files", "Package Java applications into .jar, .war, .ear files"})
	detectUnwantedFiles(data, &report, files, ".cs", "C# source code file", []string{"Do not upload C# source code. They will not be scanned", "Veracode requires the .NET application to be compiled"})
	detectUnwantedFiles(data, &report, files, ".sln", ".NET solution file", []string{"Do not upload C# source code. They will not be scanned", "Veracode requires the .NET application to be compiled"})
	detectUnwantedFiles(data, &report, files, ".csproj", "C# project file", []string{"Do not upload C# source code. They will not be scanned", "Veracode requires the .NET application to be compiled"})
	detectUnwantedFiles(data, &report, files, ".c", "C source code file", []string{"Do not upload C source code. They will not be scanned", "Veracode requires the application to be compiled with debug symbols"})
	detectUnwantedFiles(data, &report, files, ".cpp", "C++ source code file", []string{"Do not upload C++ source code. They will not be scanned", "Veracode requires the application to be compiled with debug symbols"})
	detectUnwantedFiles(data, &report, files, ".test.dll", "test artifact", []string{"Do not upload any testing artefacts"})
	detectUnwantedFiles(data, &report, files, ".unittests.dll", "test artifact", []string{"Do not upload any testing artefacts"})
	detectUnwantedFiles(data, &report, files, ".unittest.dll", "test artifact", []string{"Do not upload any testing artefacts"})
	detectUnwantedFiles(data, &report, files, ".coffee", "CoffeeScript file", []string{"CoffeeScript source code files will not be scanned", "Review the JavaScript/TypeScript packaging cheatsheet: https://nhinv11.github.io/#/JavaScript%20/%20TypeScript", "Consider using the unofficial JavaScript/TypeScript packaging tool: https://github.com/fw10/veracode-javascript-packager"})

	if report.Len() > 0 {
		printTitle("Files Uploaded")
		colorPrintf(report.String() + "\n")
	}
}

func detectSensitiveFiles(data Data, report *strings.Builder, files []string) {
	var foundFiles []string

	for _, fileName := range files {
		if isFileNameInFancyList(fileName, secretFileNames) {
			if !isStringInStringArray(fileName, foundFiles) {
				foundFiles = append(foundFiles, fileName)
			}
		}
	}

	if len(foundFiles) == 0 {
		return
	}

	report.WriteString(formatErrorStringFormat(
		"%d sensitive file%s were found: %s\n",
		len(foundFiles),
		pluralise(len(foundFiles)),
		top5StringList(foundFiles)))

	data.makeRecommendation("Do not upload any secrets, certificates or key files")
}

func detectTestArtefacts(data Data, report *strings.Builder, files []string) {
	var foundFiles []string

	for _, fileName := range files {
		if isFileNameInFancyList(fileName, testFileNames) {
			if !isStringInStringArray(fileName, foundFiles) {
				foundFiles = append(foundFiles, fileName)
			}
		}
	}

	if len(foundFiles) == 0 {
		return
	}

	report.WriteString(formatErrorStringFormat(
		"%d test artefact%s were found: %s\n",
		len(foundFiles),
		pluralise(len(foundFiles)),
		top5StringList(foundFiles)))

	data.makeRecommendation("Do not upload any testing artefacts")
}

func detectNodeModules(data Data, report *strings.Builder) {
	var foundFiles []string

	for _, file := range data.PrescanFileList.Files {
		if strings.Contains(strings.ToLower(file.Name), "_nodemodule_") {
			if !isStringInStringArray(file.Name, foundFiles) {
				foundFiles = append(foundFiles, file.Name)
			}
		}
	}

	for _, module := range data.DetailedReport.StaticAnalysis.Modules {
		if strings.Contains(strings.ToLower(module.Name), "_nodemodule_") {
			if !isStringInStringArray(module.Name, foundFiles) {
				foundFiles = append(foundFiles, module.Name)
			}
		}
	}

	if len(foundFiles) == 0 {
		return
	}

	data.makeRecommendation("Do not upload the node_modules folder as Veracode does not scan this directory")
	data.makeRecommendation("Review the JavaScript/TypeScript packaging cheatsheet: https://nhinv11.github.io/#/JavaScript%20/%20TypeScript")
	data.makeRecommendation("Consider using the unofficial JavaScript/TypeScript packaging tool: https://github.com/fw10/veracode-javascript-packager")

	report.WriteString(formatWarningString("One or more node_modules folders were detected\n"))
}

func detectGit(data Data, report *strings.Builder, files []string) {
	foundFiles := false

	for _, fileName := range files {
		if strings.EqualFold(fileName, "fsmonitor-watchman.sample") || strings.EqualFold(fileName, "FETCH_HEAD") {
			foundFiles = true
		}
	}

	if !foundFiles {
		return
	}

	report.WriteString(formatErrorString("A git repository was found\n"))

	data.makeRecommendation("Do not upload .git folders")
}

func detectRoslyn(data Data, report *strings.Builder, files []string) {
	foundFiles := false

	for _, fileName := range files {
		if strings.EqualFold(fileName, "csc.exe") {
			foundFiles = true
		}
	}

	if !foundFiles {
		return
	}

	report.WriteString(formatWarningString("The .NET Roslyn compiler was found\n"))

	data.makeRecommendation("Review the .NET packaging cheatsheet: https://nhinv11.github.io/#/.NET")
}

func detectUnwantedFiles(data Data, report *strings.Builder, files []string, suffix, name string, recommendations []string) {
	var foundFiles []string

	for _, fileName := range files {
		if strings.HasSuffix(strings.ToLower(fileName), suffix) && !isStringInStringArray(fileName, foundFiles) {
			if !isStringInStringArray(fileName, foundFiles) {
				foundFiles = append(foundFiles, fileName)
			}
		}
	}

	if len(foundFiles) == 0 {
		return
	}

	for _, recommendation := range recommendations {
		data.makeRecommendation(recommendation)
	}

	report.WriteString(formatErrorStringFormat(
		"%d %s%s: %s\n",
		len(foundFiles),
		name,
		pluralise(len(foundFiles)),
		top5StringList(foundFiles)))
}

func (data Data) reportDuplicateFiles() {
	var warningReport strings.Builder
	var errorReport strings.Builder
	var processedFiles []string

	for _, file := range data.PrescanFileList.Files {
		if isStringInStringArray(file.Name, processedFiles) {
			continue
		}

		if file.IsIgnored || file.IsThirdParty {
			continue
		}

		md5s := []string{file.MD5}
		var count = 0

		for _, otherFile := range data.PrescanFileList.Files {
			if file.Name == otherFile.Name {
				count++
				if !isStringInStringArray(otherFile.MD5, md5s) {
					md5s = append(md5s, otherFile.MD5)
				}
			}
		}

		if len(md5s) > 1 {
			if count == len(md5s) {
				warningReport.WriteString(formatWarningStringFormat(
					"%d duplicate occurance%s of \"%s\"\n",
					count,
					pluralise(count),
					file.Name))
			} else {
				errorReport.WriteString(formatErrorStringFormat(
					"%d duplicate occurance%s of \"%s\" with %d different MD5 hashes\n",
					count,
					pluralise(count),
					file.Name,
					len(md5s)))

				if strings.HasSuffix(strings.ToLower(file.Name), ".dll") || strings.HasSuffix(strings.ToLower(file.Name), ".exe") {
					data.makeRecommendation("Be careful not to upload duplicate file names with different contents as this can lead to indeterminate scan results")
				}
			}
		}

		processedFiles = append(processedFiles, file.Name)
	}

	if warningReport.Len() > 0 || errorReport.Len() > 0 {
		printTitle("Duplicate Files")
		colorPrintf(errorReport.String() + warningReport.String() + "\n")

		data.makeRecommendation("Do not upload duplicate filenames")
	}
}
