package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type PrescanModuleList struct {
	XMLName   xml.Name        `xml:"prescanresults"`
	Modules   []PrescanModule `xml:"module"`
	TotalSize int
}

type PrescanModule struct {
	XMLName        xml.Name             `xml:"module"`
	ID             int                  `xml:"id,attr"`
	Name           string               `xml:"name,attr"`
	Status         string               `xml:"status,attr"`
	Platform       string               `xml:"platform,attr"`
	Size           string               `xml:"size,attr"`
	MD5            string               `xml:"checksum,attr"`
	HasFatalErrors bool                 `xml:"has_fatal_errors,attr"`
	IsDependency   bool                 `xml:"is_dependency,attr"`
	Issues         []PrescanModuleIssue `xml:"issue"`
	IsIgnored      bool
	IsThirdParty   bool
	SizeBytes      int
}

type PrescanModuleIssue struct {
	XMLName xml.Name `xml:"issue"`
	Details string   `xml:"details,attr"`
}

func (api API) getPrescanModuleList(appId, buildId int) PrescanModuleList {
	var url = fmt.Sprintf("https://analysiscenter.veracode.com/api/5.0/getprescanresults.do?app_id=%d&build_id=%d", appId, buildId)
	response := api.makeApiRequest(url, http.MethodGet)

	moduleList := PrescanModuleList{}
	xml.Unmarshal(response, &moduleList)

	moduleList.TotalSize = 0

	for index, module := range moduleList.Modules {
		moduleList.Modules[index].IsIgnored = isFileNameInFancyList(module.Name, fileExtensionsToIgnore)
		moduleList.Modules[index].IsThirdParty = isFileNameInFancyList(module.Name, thirdPartyModules)
		moduleList.Modules[index].SizeBytes = calculateModuleSize(module.Size)
		moduleList.TotalSize += moduleList.Modules[index].SizeBytes
	}

	// Sort modules by name for consistency
	sort.Slice(moduleList.Modules, func(i, j int) bool {
		return moduleList.Modules[i].Name < moduleList.Modules[j].Name
	})

	return moduleList
}

func calculateModuleSize(size string) int {
	var totalModuleSize = 0
	totalModuleSize += convertSize(size, "GB", 1e+9)
	totalModuleSize += convertSize(size, "MB", 1e+6)
	totalModuleSize += convertSize(size, "KB", 1000)
	return totalModuleSize
}

func convertSize(size, measurement string, multiplier int) int {
	if !strings.HasSuffix(size, measurement) {
		return 0
	}

	formattedSize := strings.TrimSuffix(size, measurement)
	sizeInt, err := strconv.Atoi(formattedSize)

	if err != nil {
		panic(err)
	}

	return sizeInt * multiplier

}

func (moduleList PrescanModuleList) getFromName(moduleName string) PrescanModule {
	for _, moduleFromlist := range moduleList.Modules {
		if moduleFromlist.Name == moduleName {
			return moduleFromlist
		}
	}

	return PrescanModule{}
}

func (module PrescanModule) getFatalReason() string {
	for _, issue := range strings.Split(module.Status, ",") {
		if strings.HasPrefix(issue, "(Fatal)") {
			return strings.Replace(strings.Replace(issue, "(Fatal)", "", 1), " - 1 File", "", 1)
		}
	}

	return ""
}
