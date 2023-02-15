package main

import "strings"

var fileExtensionsToIgnore = []string{
	".cs",
	".sln",
	".asax",
	".asmx",
	".aspx",
	".manifest",
	".config",
	".xml",
	".properties",
	".md",
	".less",
}

func shouldFileNameBeIgnored(fileName string) bool {
	for _, extension := range fileExtensionsToIgnore {
		if strings.HasSuffix(strings.ToLower(fileName), extension) {
			return true
		}
	}

	return false
}
