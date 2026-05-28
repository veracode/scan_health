package report

import (
	"encoding/json"
	"fmt"
	"github.com/veracode/scan_health/v2/utils"
	"os"
	"path/filepath"
)

func (r *Report) renderToJson(filePath string) {
	formattedJson, err := json.MarshalIndent(r, "", "    ")

	if err != nil {
		utils.ErrorAndExit("Could not render to JSON", err)
	}

	if filePath != "" {
		cleanPath := filepath.Clean(filePath)
		if filepath.Ext(cleanPath) != ".json" {
			utils.ErrorAndExit("Invalid file path: output file must have a .json extension", nil)
		}
		if _, err := os.Stat(filepath.Dir(cleanPath)); os.IsNotExist(err) {
			utils.ErrorAndExit("Invalid file path: parent directory does not exist", nil)
		}
		if err := os.WriteFile(cleanPath, formattedJson, 0600); err != nil {
			utils.ErrorAndExit("Could not save JSON file", err)
		}
		return
	}

	fmt.Println(string(formattedJson))
}
