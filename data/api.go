package data

import (
	"fmt"
	"github.com/veracode/scan_health/v2/utils"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/antfie/veracode-go-hmac-authentication/hmac"
)

type API struct {
	Id            string
	Key           string
	Region        string
	AppVersion    string
	EnableCaching bool
	Profile       string
	DebugMode     bool
}

func (api API) makeApiRequest(apiUrl, httpMethod string) []byte {
	baseUrl := utils.ParseBaseUrlFromRegion(api.Region)
	fullUrl := baseUrl + apiUrl

	if api.EnableCaching {
		cachedResponse := getCachedResponse(fullUrl)

		if cachedResponse != nil {
			return cachedResponse
		}
	}

	parsedUrl, err := url.Parse(fullUrl)

	if err != nil {
		utils.ErrorAndExit("Invalid API URL", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(httpMethod, parsedUrl.String(), nil)

	if err != nil {
		utils.ErrorAndExit("Could not create API request", err)
	}

	authorizationHeader, err := hmac.CalculateAuthorizationHeader(parsedUrl, httpMethod, api.Id, api.Key)

	if err != nil {
		utils.ErrorAndExit("Could not calculate the authorization header", err)
	}

	req.Header.Add("Authorization", authorizationHeader)
	req.Header.Add("User-Agent", fmt.Sprintf("ScanHealth/%s", api.AppVersion))

	start := time.Now()

	resp, err := client.Do(req)

	if api.DebugMode {
		fmt.Printf("Request to <%s> took %.2f seconds and returned status %d\n", parsedUrl.String(), time.Since(start).Seconds(), resp.StatusCode)
	}

	if err != nil {
		utils.ErrorAndExit("There was a problem communicating with the API. Please check your connectivity and the service status page at https://status.veracode.com", err)
	}

	body, err := io.ReadAll(resp.Body)

	if api.DebugMode {
		fmt.Printf("Request to <%s> returned body %s\n", parsedUrl.String(), string(body))
	}

	if err != nil {
		utils.ErrorAndExit("There was a problem processing the API response. Please check your connectivity and the service status page at https://status.veracode.com", err)
	}

	profileFlagHint := "Do you need to specify \"-profile xyz\"?."

	if api.Profile != "default" {
		profileFlagHint = fmt.Sprintf("Is \"%s\" the right profile name to use for this region?.", api.Profile)
	}

	if resp.StatusCode == 401 {
		if strings.HasSuffix(parsedUrl.Path, "getmaintenancescheduleinfo.do") {
			utils.ErrorAndExit(fmt.Sprintf("There was a problem with your credentials. %s Please check your credentials are valid for this Veracode region. For help contact your Veracode administrator.", profileFlagHint), nil)
		} else {
			utils.ErrorAndExit(fmt.Sprintf("You are not authorized to perform this action. %s Please check you have the \"Results API\" user role set. For help contact your Veracode administrator and refer to https://docs.veracode.com/r/c_API_roles_details", profileFlagHint), nil)
		}
	}

	if resp.StatusCode == 403 {
		utils.ErrorAndExit(fmt.Sprintf("This request was forbidden. %s Ensure you can view these scans within the Veracode Platform. For help contact your Veracode administrator and refer to https://docs.veracode.com/r/c_API_roles_details", profileFlagHint), nil)
	}

	if resp.StatusCode != http.StatusOK {
		utils.ErrorAndExit(fmt.Sprintf("API responded with status of %s", resp.Status), nil)
	}

	// Cache if enabled and unless there is some indicator of an error
	if !strings.Contains(string(body), "</error>") && api.EnableCaching {
		cacheResponse(fullUrl, body)
	}

	return body
}

func (api API) AssertCredentialsWork() {
	api.makeApiRequest("/api/3.0/getmaintenancescheduleinfo.do", http.MethodGet)
}
