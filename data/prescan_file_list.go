package data

import (
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"sort"

	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
)

type prescanFileList struct {
	XMLName xml.Name       `xml:"filelist"`
	Files   []fileListFile `xml:"file"`
}

type fileListFile struct {
	XMLName xml.Name `xml:"file"`
	Id      int      `xml:"file_id,attr"`
	Name    string   `xml:"file_name,attr"`
	Status  string   `xml:"file_status,attr"`
	MD5     string   `xml:"file_md5,attr"`
}

func (api API) populatePrescanFileList(r *report.Report) {
	var url = fmt.Sprintf("/api/5.0/getfilelist.do?app_id=%d&build_id=%d", r.Scan.ApplicationId, r.Scan.BuildId)
	response := api.makeApiRequest(url, http.MethodGet)

	fileList := prescanFileList{}
	err := xml.Unmarshal(response, &fileList)

	if err != nil {
		utils.ErrorAndExit("Could not parse getfilelist.do API response", err)
	}

	// Sort files by name for consistency
	sort.Slice(fileList.Files, func(i, j int) bool {
		return fileList.Files[i].Name < fileList.Files[j].Name
	})

	for _, file := range fileList.Files {
		r.UploadedFiles = append(
			r.UploadedFiles,
			report.UploadedFile{
				Id:     file.Id,
				Name:   html.UnescapeString(file.Name),
				Status: html.UnescapeString(file.Status),
				MD5:    file.MD5,
				Source: "prescan_file_list",
			},
		)
	}
}
