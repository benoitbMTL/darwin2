package handlers

import (
	"bytes"
	"crypto/tls"
	"darwin2/config"
	"darwin2/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

// Structs to parse JSON responses
type DomainInfoResponse struct {
	Results []struct {
		DbId       int64  `json:"db_id"`
		DomainName string `json:"domain_name"`
	} `json:"results"`
}

// HandleResetMachineLearning resets the machine learning for a domain
func HandleResetMachineLearning(c echo.Context) error {
	// Initial configuration debug logs
	host := config.CurrentConfig.FWBMGTIP
	port := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()
	policyName := config.CurrentConfig.MLPOLICY
	fmt.Printf("Attempting to reset Machine Learning for policy: %s\n", policyName)

	// Create a custom HTTP client with disabled SSL verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Constructing and logging the URL for getting domain info
	domainInfoURL := fmt.Sprintf("https://%s:%s/api/v2.0/machine_learning/policy.getdomaininfo?policy_name=%s", host, port, url.QueryEscape(policyName))
	fmt.Printf("Fetching domain info from URL: %s\n", domainInfoURL)

	// Create and send the request for domain info
	req, err := http.NewRequest("GET", domainInfoURL, nil)
	if err != nil {
		errorMsg := fmt.Sprintf("Error creating request for domain info: %v", err)
		fmt.Println(errorMsg)
		return echo.NewHTTPError(http.StatusInternalServerError, errorMsg)
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		errorMsg := fmt.Sprintf("Error sending request for domain info: %v", err)
		fmt.Println(errorMsg)
		return echo.NewHTTPError(http.StatusInternalServerError, errorMsg)
	}
	defer resp.Body.Close()

	fmt.Printf("Received response with status code: %d for domain info request\n", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return err
	}

	// Parsing the domain info response
	var domainInfo DomainInfoResponse
	err = json.Unmarshal(body, &domainInfo)
	if err != nil {
		fmt.Printf("Error parsing domain info JSON: %v\n", err)
		return err
	}

	// If no domain information found
	if len(domainInfo.Results) == 0 {
		errorMsg := "No domain information found"
		fmt.Println(errorMsg)
		return echo.NewHTTPError(http.StatusNotFound, errorMsg)
	}

	dbId := domainInfo.Results[0].DbId
	domainName := domainInfo.Results[0].DomainName
	fmt.Printf("Domain %s has db_id %d. Proceeding to reset Machine Learning.\n", domainName, dbId)

	// Constructing and logging the URL for resetting ML
	resetMLURL := fmt.Sprintf("https://%s:%s/api/v2.0/machine_learning/policy.refreshdomain", host, port)
	fmt.Printf("Sending request to reset Machine Learning at URL: %s\n", resetMLURL)

	// Preparing the payload and sending the request to reset ML
	payload := fmt.Sprintf("{\"domain_id\": \"%d\", \"policy_name\": \"%s\"}", dbId, policyName)
	req, err = http.NewRequest("POST", resetMLURL, bytes.NewBufferString(payload))
	if err != nil {
		fmt.Printf("Error creating request to reset ML: %v\n", err)
		return err
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		errorMsg := fmt.Sprintf("Error sending request to reset ML: %v", err)
		fmt.Println(errorMsg)
		return echo.NewHTTPError(http.StatusInternalServerError, errorMsg)
	}
	defer resp.Body.Close()

	successMsg := fmt.Sprintf("Machine Learning for domain %s has been reset.", domainName)
	fmt.Println(successMsg)
	return c.String(http.StatusOK, successMsg)
}
