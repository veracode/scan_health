package main

import (
	"encoding/json"
	"fmt"
	"github.com/veracode/scan_health/v2/utils"
	"io"
	"net/http"

	"github.com/fatih/color"
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
var AppVersion = "0.0"

func notifyOfUpdates() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://github.com/veracode/scan_health/releases/latest", nil)

	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return
	}

	var response map[string]interface{}

	err = json.Unmarshal(body, &response)

	if err != nil {
		return
	}

	latestReleasedVersion, err := utils.StringToFloat(response["tag_name"].(string))

	if err != nil {
		return
	}

	appVersion, err := utils.StringToFloat(AppVersion)

	if err != nil {
		return
	}

	if latestReleasedVersion > appVersion {
		color.HiYellow(fmt.Sprintf("Please upgrade to the latest version of this tool (v%s) by visiting https://github.com/veracode/scan_health/releases/latest\n", response["tag_name"]))
	}
}
