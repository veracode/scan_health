package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlatformUrl(t *testing.T) {
	t.Run("Valid Platform Url", func(t *testing.T) {
		t.Parallel()
		assert.True(t, isPlatformURL("https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:41807:380748:24113946:5206587:24100796::::469963"))
	})

	t.Run("Valid Url, not the platform", func(t *testing.T) {
		t.Parallel()
		assert.False(t, isPlatformURL("https://test.com/"))
	})

	t.Run("Valid Platform Url, mixed scheme/host case", func(t *testing.T) {
		t.Parallel()
		assert.True(t, isPlatformURL("hTtps://aNalysiscenter.verAcode.cOm/auth/index.jsp#AnalyzeAppModuleList:41807:380748:24113946:5206587:24100796::::469963"))
	})

	t.Run("Valid Platform Domain, not a valid path", func(t *testing.T) {
		t.Parallel()
		assert.False(t, isPlatformURL("https://analysiscenter.veracode.com/help"))
	})

	t.Run("Check Commercial Region From Url", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, ParseRegionFromUrl("https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:41807:380748:24113946:5206587:24100796::::469963"), "commercial")
	})

	t.Run("Check EU Region From Url", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, ParseRegionFromUrl("https://analysiscenter.veracode.eu/auth/index.jsp#AnalyzeAppModuleList:41807:380748:24113946:5206587:24100796::::469963"), "european")
	})

	t.Run("Check US Region From Url", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, ParseRegionFromUrl("https://analysiscenter.veracode.us/auth/index.jsp#AnalyzeAppModuleList:41807:380748:24113946:5206587:24100796::::469963"), "us")
	})

	t.Run("Get Base URL for Commercial", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, ParseBaseUrlFromRegion("commercial"), "https://analysiscenter.veracode.com")
	})

	t.Run("Get Base URL defaults to Commercial", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, ParseBaseUrlFromRegion("region_does_not_exist"), "https://analysiscenter.veracode.com")
	})

	t.Run("Get Base URL for EU", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, ParseBaseUrlFromRegion("european"), "https://analysiscenter.veracode.eu")
	})

	t.Run("Get Build ID from Valid Build ID", func(t *testing.T) {
		t.Parallel()
		buildId, err := ParseBuildIdFromScanInformation("12345678")
		assert.Equal(t, buildId, 12345678)
		assert.Nil(t, err)

	})

	t.Run("Get Build ID from Valid URL", func(t *testing.T) {
		t.Parallel()
		buildId, err := ParseBuildIdFromScanInformation("https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:41807:380748:12345678:5206587:24100796::::469963")
		assert.Equal(t, buildId, 12345678)
		assert.Nil(t, err)
	})

	t.Run("Get Error From Negative Build ID", func(t *testing.T) {
		t.Parallel()
		buildId, err := ParseBuildIdFromScanInformation("-40")
		assert.Equal(t, buildId, -1)
		assert.Contains(t, err.Error(), "build ID must be positive")
	})

	t.Run("Get Error From Build ID in invalid URL", func(t *testing.T) {
		t.Parallel()
		buildId, err := ParseBuildIdFromScanInformation("https://analysisceasdasdnter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:41807:380748:12345678:5206587:24100796::::469963")
		assert.Equal(t, buildId, -1)
		assert.Contains(t, err.Error(), "not a valid or supported Veracode Platform URL")
	})

	t.Run("Checking A Valid Region", func(t *testing.T) {
		t.Parallel()
		validRegion := IsValidRegion("european")
		assert.True(t, validRegion)
	})

	t.Run("Checking An Invalid Region", func(t *testing.T) {
		t.Parallel()
		validRegion := IsValidRegion("apac")
		assert.False(t, validRegion)
	})

	t.Run("Should parse applciation ID from platform URL", func(t *testing.T) {
		t.Parallel()
		applicationId, err := ParseApplicationIdFromPlatformUrl("https://analysiscenter.veracode.com/auth/index.jsp#ReviewResultsStaticFlaws:29692:771572:33659037:33628505:33644155:::2018345")

		assert.Nil(t, err)
		assert.Equal(t, 771572, applicationId)
	})
}
