package handlers

import (
	"bytes"
	"crypto/tls"
	"darwin2/config"
	"darwin2/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

// Structs to parse JSON responses
type DomainInfoResponse struct {
	Results []struct {
		DbId       string `json:"db_id"`
		DomainName string `json:"domain_name"`
	} `json:"results"`
}

// HandleResetMachineLearning resets the machine learning for a domain
func HandleResetMachineLearning(c echo.Context) error {
	host := config.CurrentConfig.FWBMGTIP
	port := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()
	policyName := config.CurrentConfig.MLPOLICY

	fmt.Printf("Host: %s\n", host)
	fmt.Printf("Port: %s\n", port)
	fmt.Printf("Token: %s\n", token)
	fmt.Printf("Policy Name: %s\n", policyName)

	// Create a custom HTTP client with disabled SSL verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Construct the URL for getting domain info
	domainInfoURL := fmt.Sprintf("https://%s:%s/api/v2.0/machine_learning/policy.getdomaininfo?policy_name=%s", host, port, url.QueryEscape(policyName))

	// Create the request for domain info
	req, err := http.NewRequest("GET", domainInfoURL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("accept", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parse the response
	var domainInfo DomainInfoResponse
	err = json.Unmarshal(body, &domainInfo)
	if err != nil {
		return err
	}

	if len(domainInfo.Results) == 0 {
		return fmt.Errorf("no domain information found")
	}

	dbId := domainInfo.Results[0].DbId
	domainName := domainInfo.Results[0].DomainName

	fmt.Printf("Domain %s has db_id %s\n", domainName, dbId)
	fmt.Printf("Resetting Machine Learning for Domain %s\n", domainName)

	// Construct the URL for resetting ML
	resetMLURL := fmt.Sprintf("https://%s:%s/api/v2.0/machine_learning/policy.refreshdomain", host, port)

	// Prepare the payload
	payload := fmt.Sprintf("{\"domain_id\": \"%s\", \"policy_name\": \"%s\"}", dbId, policyName)
	req, err = http.NewRequest("POST", resetMLURL, bytes.NewBufferString(payload))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Send the request to reset ML
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Machine Learning for domain " + domainName + " has been reset.")

	return c.String(http.StatusOK, "Machine Learning for domain "+domainName+" has been reset.")
}
