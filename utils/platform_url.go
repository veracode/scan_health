package utils

import (
	"fmt"
	"strconv"
	"strings"
)

var supportedPages = []string{
	"ReviewResultsStaticFlaws",
	"ReviewResultsAllFlaws",
	"AnalyzeAppModuleList",
	"StaticOverview",
	"AnalyzeAppSourceFiles",
	"ViewReportsResultSummary",
	"ViewReportsDetailedReport"}

func PlatformUrlInvalid(url string) {
	ErrorAndExit(fmt.Sprintf("%s is not a valid or supported Veracode Platform URL.\nThis tool requires a URL to one of the following Veracode Platform pages: %s",
		url,
		strings.Join(supportedPages, ", ")), nil)
}

func IsPlatformURL(url string) bool {
	return strings.HasPrefix(url, "https://analysiscenter.veracode.com/auth/index.jsp") ||
		strings.HasPrefix(url, "https://analysiscenter.veracode.us/auth/index.jsp") ||
		strings.HasPrefix(url, "https://analysiscenter.veracode.eu/auth/index.jsp")
}

func IsParseableURL(urlFragment string) bool {
	for _, page := range supportedPages {
		if strings.HasPrefix(urlFragment, page) {
			return true
		}
	}
	return false
}

func ParseRegionFromUrl(url string) string {
	if strings.HasPrefix(url, "https://analysiscenter.veracode.us") {
		return "us"
	}

	if strings.HasPrefix(url, "https://analysiscenter.veracode.eu") {
		return "european"
	}

	return "commercial"
}

func ParseBaseUrlFromRegion(region string) string {
	if region == "us" {
		return "https://analysiscenter.veracode.us"
	}

	if region == "european" {
		return "https://analysiscenter.veracode.eu"
	}

	return "https://analysiscenter.veracode.com"
}

func ParseBuildIdFromPlatformUrl(urlOrBuildId string) int {
	buildId, err := strconv.Atoi(urlOrBuildId)

	if err == nil {
		return buildId
	}

	if !IsPlatformURL(urlOrBuildId) {
		PlatformUrlInvalid(urlOrBuildId)
	}

	var urlFragment = strings.Split(urlOrBuildId, "#")[1]

	if IsParseableURL(urlFragment) {
		buildId, err := strconv.Atoi(strings.Split(urlFragment, ":")[3])

		if err != nil {
			PlatformUrlInvalid(urlOrBuildId)
		}

		return buildId

	}

	PlatformUrlInvalid(urlOrBuildId)
	return -1
}
