package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var supportedPages = []string{
	"ReviewResultsStaticFlaws",
	"ReviewResultsAllFlaws",
	"AnalyzeAppModuleList",
	"StaticOverview",
	"AnalyzeAppSourceFiles",
	"ViewReportsResultSummary",
	"ViewReportsDetailedReport"}

func platformUrlInvalid(url string) {
	color.HiRed(
		fmt.Sprintf("%s is not a valid or supported Veracode Platform URL.\nThis tool requires a URL to one of the following Veracode Platform pages: %s",
			url,
			strings.Join(supportedPages, ", ")))
	os.Exit(1)
}

func isPlatformURL(url string) bool {
	return strings.HasPrefix(url, "https://analysiscenter.veracode.com/auth/index.jsp") ||
		strings.HasPrefix(url, "https://analysiscenter.veracode.us/auth/index.jsp") ||
		strings.HasPrefix(url, "https://analysiscenter.veracode.eu/auth/index.jsp")
}

func isParseableURL(urlFragment string) bool {
	for _, page := range supportedPages {
		if strings.HasPrefix(urlFragment, page) {
			return true
		}
	}
	return false
}

func parseRegionFromUrl(url string) string {
	if strings.HasPrefix(url, "https://analysiscenter.veracode.us") {
		return "us"
	}

	if strings.HasPrefix(url, "https://analysiscenter.veracode.eu") {
		return "european"
	}

	return "commercial"
}

func parseBaseUrlFromRegion(region string) string {
	if region == "us" {
		return "https://analysiscenter.veracode.us"
	}

	if region == "european" {
		return "https://analysiscenter.veracode.eu"
	}

	return "https://analysiscenter.veracode.com"
}

func parseAccountIdFromPlatformUrl(urlOrAccountId string) int {
	accountId, err := strconv.Atoi(urlOrAccountId)

	if err == nil {
		return accountId
	}

	if !isPlatformURL(urlOrAccountId) {
		platformUrlInvalid(urlOrAccountId)
	}

	var urlFragment = strings.Split(urlOrAccountId, "#")[1]

	if isParseableURL(urlFragment) {
		accountId, err := strconv.Atoi(strings.Split(urlFragment, ":")[1])

		if err != nil {
			platformUrlInvalid(urlOrAccountId)
		}

		return accountId

	}

	platformUrlInvalid(urlOrAccountId)
	return -1
}

func parseAppIdFromPlatformUrl(urlOrBuildId string) int {
	_, err := strconv.Atoi(urlOrBuildId)

	if err == nil {
		// This is a build ID, not an app ID. We cannot resolve the app ID.
		return -1
	}

	if !isPlatformURL(urlOrBuildId) {
		platformUrlInvalid(urlOrBuildId)
	}

	var urlFragment = strings.Split(urlOrBuildId, "#")[1]

	if isParseableURL(urlFragment) {
		appId, err := strconv.Atoi(strings.Split(urlFragment, ":")[2])

		if err != nil {
			platformUrlInvalid(urlOrBuildId)
		}

		return appId

	}

	platformUrlInvalid(urlOrBuildId)
	return -1
}

func parseBuildIdFromPlatformUrl(urlOrBuildId string) int {
	buildId, err := strconv.Atoi(urlOrBuildId)

	if err == nil {
		return buildId
	}

	if !isPlatformURL(urlOrBuildId) {
		platformUrlInvalid(urlOrBuildId)
	}

	var urlFragment = strings.Split(urlOrBuildId, "#")[1]

	if isParseableURL(urlFragment) {
		buildId, err := strconv.Atoi(strings.Split(urlFragment, ":")[3])

		if err != nil {
			platformUrlInvalid(urlOrBuildId)
		}

		return buildId

	}

	platformUrlInvalid(urlOrBuildId)
	return -1
}
