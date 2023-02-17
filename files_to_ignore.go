package main

import "strings"

var fileExtensionsToIgnore = []string{
	"*.asax",
	"*.asmx",
	"*.aspx",
	"*.config",
	"*.cs",
	"*.less",
	"*.manifest",
	"*.md",
	"*.png",
	"*.properties",
	"*.sln",
	"*.txt",
	"*.xml",
	"AUTHORS",
	"LICENCE",
	"Makefile",
}

func shouldFileNameBeIgnored(fileName string) bool {
	for _, extension := range fileExtensionsToIgnore {
		if strings.HasSuffix(strings.ToLower(fileName), extension) {
			return true
		}
	}

	return false
}
