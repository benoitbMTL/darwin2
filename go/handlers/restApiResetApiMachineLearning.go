package handlers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"darwin2/config"
	"darwin2/utils"

	"github.com/labstack/echo/v4"
)

// HandleResetApiMachineLearning resets the API machine learning
func HandleResetApiMachineLearning(c echo.Context) error {
	host := config.CurrentConfig.FWBMGTIP
	port := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()

	// Shared custom HTTP client with disabled SSL verification
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Fetch policy rules
	getPolicyRuleURL := fmt.Sprintf("https://%s:%s/api/v2.0/machine_learning/api_learning_policy.get_policy_rule", host, port)
	req, err := http.NewRequest("GET", getPolicyRuleURL, nil)
	if err != nil {
		fmt.Printf("Error creating GET request: %v\n", err)
		return err
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending GET request: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return err
	}

	var response map[string][]map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Error parsing JSON response: %v\n", err)
		return err
	}

	var resetDomains []string

	// Reset machine learning for each domain
	for _, result := range response["results"] {
		rules := result["rule"].([]interface{})
		for _, rule := range rules {
			ruleMap := rule.(map[string]interface{})
			ruleID := ruleMap["id"]
			domainName := ruleMap["domain-name"].(string)

			resetMLURL := fmt.Sprintf("https://%s:%s/api/v2.0/machine_learning/api_learning_policy.refreshdomain?rule_id=%v", host, port, ruleID)
			req, err := http.NewRequest("POST", resetMLURL, nil)
			if err != nil {
				fmt.Printf("Error creating POST request for rule ID %v: %v\n", ruleID, err)
				continue
			}
			req.Header.Add("Authorization", token)
			req.Header.Add("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error sending POST request for rule ID %v: %v\n", ruleID, err)
				continue
			}
			resp.Body.Close()

			resetDomains = append(resetDomains, domainName)
			fmt.Printf("Machine Learning for domain %s has been reset successfully.\n", domainName)
		}
	}

	// Return a string with all domains reset
	return c.String(http.StatusOK, fmt.Sprintf("ML for %s has been reset.", url.QueryEscape(strings.Join(resetDomains, ", "))))
}
