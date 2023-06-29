package report

import (
	"encoding/json"
	"fmt"
	"github.com/antfie/scan_health/v2/utils"
)

func (r *Report) renderToJson() {
	formattedJson, err := json.MarshalIndent(r, "", "    ")

	if err != nil {
		utils.ErrorAndExit("Could not render to JSON", err)
	}

	fmt.Println(string(formattedJson))
}
