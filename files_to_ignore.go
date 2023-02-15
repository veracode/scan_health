package main

import "strings"

var fileExtensionsToIgnore = []string{
	".asax",
	".asmx",
	".aspx",
	".config",
	".cs",
	".less",
	".manifest",
	".md",
	".properties",
	".sln",
	".txt",
	".xml",
}

func shouldFileNameBeIgnored(fileName string) bool {
	for _, extension := range fileExtensionsToIgnore {
		if strings.HasSuffix(strings.ToLower(fileName), extension) {
			return true
		}
	}

	return false
}
