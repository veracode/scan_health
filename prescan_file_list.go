package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
)

type PrescanFileList struct {
	XMLName xml.Name      `xml:"filelist"`
	Files   []PrescanFile `xml:"file"`
}

type PrescanFile struct {
	XMLName      xml.Name `xml:"file"`
	ID           int      `xml:"file_id,attr"`
	Name         string   `xml:"file_name,attr"`
	Status       string   `xml:"file_status,attr"`
	MD5          string   `xml:"file_md5,attr"`
	IsIgnored    bool
	IsThirdParty bool
}

func (api API) getPrescanFileList(appId, buildId int) PrescanFileList {
	var url = fmt.Sprintf("https://analysiscenter.veracode.com/api/5.0/getfilelist.do?app_id=%d&build_id=%d", appId, buildId)
	response := api.makeApiRequest(url, http.MethodGet)

	fileList := PrescanFileList{}
	xml.Unmarshal(response, &fileList)

	for index, file := range fileList.Files {
		fileList.Files[index].IsIgnored = isFileNameInFancyList(file.Name, fileExtensionsToIgnore)
		fileList.Files[index].IsThirdParty = isFileNameInFancyList(file.Name, thirdPartyModules)
	}

	// Sort files by name for consistency
	sort.Slice(fileList.Files, func(i, j int) bool {
		return fileList.Files[i].Name < fileList.Files[j].Name
	})

	return fileList
}

func (fileList PrescanFileList) getFromName(moduleName string) PrescanFile {
	for _, file := range fileList.Files {
		if file.Name == moduleName {
			return file
		}
	}

	return PrescanFile{}
}
