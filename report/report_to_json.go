package report

import (
	"encoding/json"
	"fmt"
	"github.com/veracode/scan_health/v2/utils"
	"os"
)

func (r *Report) renderToJson(filePath string) {
	formattedJson, err := json.MarshalIndent(r, "", "    ")

	if err != nil {
		utils.ErrorAndExit("Could not render to JSON", err)
	}

	if filePath != "" {
		if err := os.WriteFile(filePath, formattedJson, 0600); err != nil {
			utils.ErrorAndExit("Could not save JSON file", err)
		}
		return
	}

	fmt.Println(string(formattedJson))
}
