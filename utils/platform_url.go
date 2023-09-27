package utils

import (
	"errors"
	"fmt"
	"net/url"
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
	"ViewReportsDetailedReport",
	"ReviewResultsSCA"}

type Region struct {
	ID  string
	URL string
}

var regions = []Region{
	{ID: "european", URL: "https://analysiscenter.veracode.eu"},
	{ID: "us", URL: "https://analysiscenter.veracode.us"},
	{ID: "commercial", URL: "https://analysiscenter.veracode.com"},
}

var defaultRegion = regions[2]

func reportPlatformUrlInvalid(url string) string {
	return fmt.Sprintf("%s is not a valid or supported Veracode Platform URL.\nThis tool requires a URL to one of the following Veracode Platform pages: %s",
		url,
		strings.Join(supportedPages, ", "))
}

func IsPlatformURL(platformUrl string) bool {

	// We need to normalize, as schemes and hosts may be
	// variable casing (as per RFC3986)
	parsedInputURL, err := url.Parse(strings.ToLower(platformUrl))
	if err != nil {
		return false
	}

	baseInputURL := fmt.Sprintf("%s://%s%s",
		strings.ToLower(parsedInputURL.Scheme),
		strings.ToLower(parsedInputURL.Host),
		parsedInputURL.Path)

	for _, region := range regions {
		if strings.HasPrefix(baseInputURL, region.URL+"/auth/index.jsp") {
			return true
		}
	}

	return false
}

func ParseRegionFromUrl(url string) string {

	for _, region := range regions {
		if strings.HasPrefix(strings.ToLower(url), region.URL) {
			return region.ID
		}
	}

	// Should we raise an error here instead?
	return defaultRegion.ID
}

func ParseBaseUrlFromRegion(region string) string {

	for _, regionData := range regions {
		if regionData.ID == region {
			return regionData.URL
		}
	}

	return defaultRegion.URL
}

func isParseableURL(urlFragment string) bool {
	for _, page := range supportedPages {
		if strings.HasPrefix(urlFragment, page) {
			return true
		}
	}
	return false
}

func ParseBuildIdFromScanInformation(urlOrBuildId string) (int, error) {
	buildId, err := strconv.Atoi(urlOrBuildId)

	if err == nil {
		if buildId <= 0 {
			return -1, errors.New("build ID must be positive")
		}

		return buildId, nil
	}

	if !IsPlatformURL(urlOrBuildId) {
		return -1, errors.New(reportPlatformUrlInvalid(urlOrBuildId))
	}

	var urlFragment = strings.Split(urlOrBuildId, "#")[1]

	if isParseableURL(urlFragment) {
		buildId, err := strconv.Atoi(strings.Split(urlFragment, ":")[3])

		if err != nil {
			return -1, errors.New(reportPlatformUrlInvalid(urlOrBuildId))
		}

		return buildId, nil

	}

	return -1, errors.New(reportPlatformUrlInvalid(urlOrBuildId))
}
