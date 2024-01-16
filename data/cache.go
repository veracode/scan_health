package data

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/veracode/scan_health/v2/utils"
	"os"
	"path"
	"path/filepath"
)

func getCachedResponse(url string) []byte {
	hash := getHash(url)
	filePath, err := filepath.Abs(path.Join("cache", hash))

	if err != nil {
		utils.ErrorAndExit("Could not resolve absolute path", err)
	}

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	data, err := os.ReadFile(filePath)

	if err != nil {
		utils.ErrorAndExit("Could not load data from cache", err)
	}

	return data
}

var cacheDirectoryInitialized = false

func cacheResponse(url string, data []byte) {
	if !cacheDirectoryInitialized {
		if _, err := os.Stat("cache"); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir("cache", 0777)

			if err != nil {
				utils.ErrorAndExit("Could not create cache path", err)
			}
		}

		cacheDirectoryInitialized = true
	}

	hash := getHash(url)
	filePath, err := filepath.Abs(path.Join("cache", hash))

	if err != nil {
		utils.ErrorAndExit("Could not resolve absolute path", err)
	}

	err = os.WriteFile(filePath, data, 0777)

	if err != nil {
		utils.ErrorAndExit("Could not cache data", err)
	}
}

func getHash(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
