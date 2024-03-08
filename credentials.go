package main

import (
	"errors"
	"fmt"
	"github.com/veracode/scan_health/v2/utils"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

func formatCredential(credential string) string {
	var parts = strings.Split(credential, "-")
	if len(parts) == 2 {
		return parts[1]
	}

	return credential
}

func getCredentials(id, key string, profile string) (string, string) {
	id = formatCredential(id)
	key = formatCredential(key)

	// First try CLI flags
	if len(id) == 32 && len(key) == 128 {
		return id, key
	}

	if len(id) > 0 && len(id) != 32 {
		utils.ErrorAndExit("Invalid value for -vid", nil)
	}

	if len(key) > 0 && len(key) != 128 {
		utils.ErrorAndExit("Invalid value for -vkey", nil)
	}

	if len(id) > 0 && len(key) == 0 || len(key) > 0 && len(id) == 0 {
		utils.ErrorAndExit("If passing Veracode API key via command line both -vid and -vkey are required", nil)
	}

	id = ""
	key = ""

	// Then try environment variables
	id = os.Getenv("VERACODE_API_KEY_ID")
	key = os.Getenv("VERACODE_API_KEY_SECRET")

	id = formatCredential(id)
	key = formatCredential(key)

	if len(id) == 32 && len(key) == 128 {
		return id, key
	}

	if len(id) > 0 && len(id) != 32 {
		utils.ErrorAndExit("Invalid value for VERACODE_API_KEY_ID", nil)
	}

	if len(key) > 0 && len(key) != 128 {
		utils.ErrorAndExit("Invalid value for VERACODE_API_KEY_SECRET", nil)
	}

	if len(id) > 0 && len(key) == 0 || len(key) > 0 && len(id) == 0 {
		utils.ErrorAndExit("If passing Veracode API key via environment variables both VERACODE_API_KEY_ID and VERACODE_API_KEY_SECRET are required", nil)
	}

	id = ""
	key = ""

	// Finally look for a Veracode credentials file
	homePath, err := os.UserHomeDir()

	if err != nil {
		utils.ErrorAndExit("Could not locate your home directory", err)

	}

	var credentialsFilePath = filepath.Join(homePath, ".veracode", "credentials")

	if _, err := os.Stat(credentialsFilePath); errors.Is(err, os.ErrNotExist) {
		utils.ErrorAndExit("Could not resolve any API credentials. Use either -vid and -vkey command line arguments, set VERACODE_API_KEY_ID and VERACODE_API_KEY_SECRET environment variables or create a Veracode credentials file. See https://docs.veracode.com/r/c_api_credentials3", err)

	}

	cfg, err := ini.Load(credentialsFilePath)
	if err != nil {
		utils.ErrorAndExit("Could not open the Veracode credentials file. See https://docs.veracode.com/r/c_api_credentials3", err)

	}

	if !cfg.HasSection(profile) {
		profileFlagHint := ""

		if profile == "default" {
			profileFlagHint = "Do you need to specify \"-profile xyz\"?. "
		}

		utils.ErrorAndExit(fmt.Sprintf("Could not find the profile [%s] within the Veracode credentials file. %sSee https://docs.veracode.com/r/c_httpie_tool#using-multiple-profiles", profile, profileFlagHint), nil)
	}

	id = cfg.Section(profile).Key("veracode_api_key_id").String()
	key = cfg.Section(profile).Key("veracode_api_key_secret").String()

	if len(id) > 0 && len(key) > 0 {
		id = formatCredential(id)
		key = formatCredential(key)

		if len(id) != 32 {
			utils.ErrorAndExit(fmt.Sprintf("Invalid value for veracode_api_key_id in file \"%s\"", credentialsFilePath), nil)

		}

		if len(key) != 128 {
			utils.ErrorAndExit(fmt.Sprintf("Invalid value for veracode_api_key_secret in file \"%s\"", credentialsFilePath), nil)

		}

		return id, key
	}

	utils.ErrorAndExit("Could not parse credentials from the Veracode credentials file. See https://docs.veracode.com/r/c_api_credentials3", nil)

	return "", ""
}
