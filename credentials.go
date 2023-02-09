package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
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
		color.HiRed("Error: Invalid value for -vid")
		os.Exit(1)
	}

	if len(key) > 0 && len(key) != 128 {
		color.HiRed("Error: Invalid value for -vkey")
		os.Exit(1)
	}

	if len(id) > 0 && len(key) == 0 || len(key) > 0 && len(id) == 0 {
		color.HiRed("Error: If passing Veracode API key via command line both -vid and -vkey are required")
		os.Exit(1)
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
		color.HiRed("Error: Invalid value for VERACODE_API_KEY_ID")
		os.Exit(1)
	}

	if len(key) > 0 && len(key) != 128 {
		color.HiRed("Error: Invalid value for VERACODE_API_KEY_SECRET")
		os.Exit(1)
	}

	if len(id) > 0 && len(key) == 0 || len(key) > 0 && len(id) == 0 {
		color.HiRed("Error: If passing Veracode API key via environment variables both VERACODE_API_KEY_ID and VERACODE_API_KEY_SECRET are required")
		os.Exit(1)
	}

	id = ""
	key = ""

	// Finally look for a Veracode credentials file
	homePath, err := os.UserHomeDir()

	if err != nil {
		color.HiRed("Error: Could not locate your home directory")
		os.Exit(1)
	}

	var credentialsFilePath = filepath.Join(homePath, ".veracode", "credentials")

	if _, err := os.Stat(credentialsFilePath); errors.Is(err, os.ErrNotExist) {
		color.HiRed("Error: Could not resolve any API credentials. Use either -vid and -vkey command line arguments, set VERACODE_API_KEY_ID and VERACODE_API_KEY_SECRET environment variables or create a Veracode credentials file. See https://docs.veracode.com/r/c_configure_api_cred_file")
		os.Exit(1)
	}

	cfg, err := ini.Load(credentialsFilePath)
	if err != nil {
		color.HiRed("Error: Could not open the Veracode credentials file. See https://docs.veracode.com/r/c_configure_api_cred_file")
		os.Exit(1)
	}

	if !cfg.HasSection(profile) {
		color.HiRed(fmt.Sprintf("Error: Could not find the profile [%s] within the Veracode credentials file. See https://docs.veracode.com/r/c_httpie_tool", profile))
		os.Exit(1)
	}

	id = cfg.Section(profile).Key("veracode_api_key_id").String()
	key = cfg.Section(profile).Key("veracode_api_key_secret").String()

	if len(id) > 0 && len(key) > 0 {
		id = formatCredential(id)
		key = formatCredential(key)

		if len(id) != 32 {
			color.HiRed("Error: Invalid value for veracode_api_key_id in file \"%s\"", credentialsFilePath)
			os.Exit(1)
		}

		if len(key) != 128 {
			color.HiRed("Error: Invalid value for veracode_api_key_secret in file \"%s\"", credentialsFilePath)
			os.Exit(1)
		}

		return id, key
	}

	color.HiRed("Error: Could not parse credentials from the Veracode credentials file. See https://docs.veracode.com/r/c_configure_api_cred_file")
	os.Exit(1)
	return "", ""
}
