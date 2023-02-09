package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fatih/color"
)

var AppVersion string = "0.0"

func notifyOfUpdates() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://github.com/antfie/scan_health/releases/latest", nil)

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

	latestReleasedVersion, err := stringToFloat(response["tag_name"].(string))

	if err != nil {
		return
	}

	appVeresion, err := stringToFloat(AppVersion)

	if err != nil {
		return
	}

	if latestReleasedVersion > appVeresion {
		color.HiYellow(fmt.Sprintf("Please upgrade to the latest version of this tool (v%s) by visiting https://github.com/antfie/scan_health/releases/latest\n", response["tag_name"]))
	}
}
