package main

import "strings"

var thirdPartyModules = []string{
	"componentspace.saml2.dll",
	"devexpress.*",
	"entityframework.dll",
	"itextsharp.dll",
	"log4net.dll",
	"microsoft.*.dll",
	"microsoft.*.pdb",
	"newrelic.*.dll",
	"newtonsoft.json.dll",
	"ninject.web.common.webhost.dll",
	"syncfusion.*",
	"system.*.dll",
}

func isThirdParty(fileName string) bool {
	formattedFileName := strings.TrimSpace(strings.ToLower(fileName))

	for _, thirdPartyModule := range thirdPartyModules {
		// There can only be one wildcard
		if strings.Count(thirdPartyModule, "*") == 1 {
			// At the start
			if strings.HasPrefix(thirdPartyModule, "*") {
				if strings.HasSuffix(formattedFileName, strings.ReplaceAll(thirdPartyModule, "*", "")) {
					return true
				}

				// At the end
			} else if strings.HasSuffix(thirdPartyModule, "*") {
				if strings.HasPrefix(formattedFileName, strings.ReplaceAll(thirdPartyModule, "*", "")) {
					return true
				}

				// Or somewhere in the middle
			} else {
				parts := strings.Split(thirdPartyModule, "*")

				if strings.HasPrefix(formattedFileName, parts[0]) && strings.HasSuffix(formattedFileName, parts[1]) {
					return true
				}
			}
		} else if strings.EqualFold(formattedFileName, thirdPartyModule) {
			return true
		}
	}

	return false
}
