package data

import (
	"fmt"
	"github.com/antfie/scan_health/v2/utils"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/antfie/veracode-go-hmac-authentication/hmac"
)

type API struct {
	Id            string
	Key           string
	Region        string
	AppVersion    string
	EnableCaching bool
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

	resp, err := client.Do(req)

	if err != nil {
		utils.ErrorAndExit("There was a problem communicating with the API. Please check your connectivity and the service status page at https://status.veracode.com", err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		utils.ErrorAndExit("There was a problem processing the API response. Please check your connectivity and the service status page at https://status.veracode.com", err)
	}

	if resp.StatusCode == 401 {
		if strings.HasSuffix(parsedUrl.Path, "getmaintenancescheduleinfo.do") {
			utils.ErrorAndExit("There was a problem with your credentials. Please check your credentials are valid for this Veracode region. For help contact your Veracode administrator.", nil)
		} else {
			utils.ErrorAndExit("You are not authorized to perform this action. Please check you have the \"Results API\" user role set. For help contact your Veracode administrator and refer to https://docs.veracode.com/r/c_API_roles_details", nil)
		}
	}

	if resp.StatusCode == 403 {
		utils.ErrorAndExit("This request was forbidden. Ensure you can view these scans within the Veracode Platform. For help contact your Veracode administrator and refer to https://docs.veracode.com/r/c_API_roles_details", nil)
	}

	if resp.StatusCode != http.StatusOK {
		utils.ErrorAndExit(fmt.Sprintf("API responded with status of %s", resp.Status), nil)
	}

	if api.EnableCaching {
		cacheResponse(fullUrl, body)
	}

	return body
}

func (api API) AssertCredentialsWork() {
	api.makeApiRequest("/api/3.0/getmaintenancescheduleinfo.do", http.MethodGet)
}
