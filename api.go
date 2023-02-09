package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/antfie/veracode-go-hmac-authentication/hmac"
	"github.com/fatih/color"
)

type API struct {
	id     string
	key    string
	region string
}

func (api API) makeApiRequest(apiUrl, httpMethod string) []byte {
	if api.region == "us" {
		apiUrl = strings.Replace(apiUrl, ".com", ".us", 1)
	} else if api.region == "eu" {
		apiUrl = strings.Replace(apiUrl, ".com", ".eu", 1)
	}

	parsedUrl, err := url.Parse(apiUrl)

	if err != nil {
		color.HiRed("Error: Invalid API URL")
		os.Exit(1)
	}

	client := &http.Client{}
	req, err := http.NewRequest(httpMethod, parsedUrl.String(), nil)

	if err != nil {
		color.HiRed("Error: Could not create API request")
		os.Exit(1)
	}

	authorizationHeader, err := hmac.CalculateAuthorizationHeader(parsedUrl, httpMethod, api.id, api.key)

	if err != nil {
		color.HiRed("Error: Could not calculate the authorization header")
		os.Exit(1)
	}

	req.Header.Add("Authorization", authorizationHeader)
	req.Header.Add("User-Agent", fmt.Sprintf("ScanHealth/%s", AppVersion))

	resp, err := client.Do(req)

	if err != nil {
		color.HiRed("Error: There was a problem communicating with the API. Please check your connectivity and the service status page at https://status.veracode.com")
		os.Exit(1)
	}

	if resp.StatusCode == 401 {
		if strings.HasSuffix(parsedUrl.Path, "getmaintenancescheduleinfo.do") {
			color.HiRed("Error: There was a problem with your credentials. Please check your credentials are valid for this Veracode region. For help contact your Veracode administrator.")
		} else {
			color.HiRed("Error: You are not authorized to perform this action. Please check you have the \"Results API\" user role set. For help contact your Veracode administrator and refer to https://docs.veracode.com/r/c_API_roles_details")
		}

		os.Exit(1)
	}

	if resp.StatusCode == 403 {
		color.HiRed("Error: This request was forbidden. Ensure you can view these scans within the Veracode Platform. For help contact your Veracode administrator and refer to https://docs.veracode.com/r/c_API_roles_details")
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		color.HiRed(fmt.Sprintf("Error: API request returned status of %s", resp.Status))
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		color.HiRed("Error: There was a problem processing the API response. Please check your connectivity and the service status page at https://status.veracode.com")
		os.Exit(1)
	}

	return body
}

func (api API) assertCredentialsWork() {
	api.makeApiRequest("https://analysiscenter.veracode.com/api/3.0/getmaintenancescheduleinfo.do", http.MethodGet)
}
