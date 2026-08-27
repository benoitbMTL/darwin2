package handlers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"darwin2/config"
	"darwin2/utils"

	"github.com/labstack/echo/v4"
)

// HandleResetApiMachineLearning resets the API machine learning
func HandleResetApiMachineLearning(c echo.Context) error {
	currentConfig := config.GetCurrentConfig()
	host := currentConfig.FWBMGTIP
	port := currentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()

	// curl opens a new connection for every command in the working test script.
	// Do the same here because FortiWeb may not handle reusing the GET connection
	// for the refresh POST requests reliably.
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
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
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return c.String(http.StatusBadGateway, fmt.Sprintf("Unable to read FortiWeb policy rules response: %v", err))
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		errorMsg := fmt.Sprintf("FortiWeb policy rules request failed with status %s: %s", resp.Status, strings.TrimSpace(string(body)))
		fmt.Println(errorMsg)
		return c.String(http.StatusBadGateway, errorMsg)
	}

	var response map[string][]map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Error parsing JSON response: %v\n", err)
		return c.String(http.StatusBadGateway, fmt.Sprintf("Invalid FortiWeb policy rules response: %v", err))
	}

	var resetDomains []string
	var resetErrors []string

	// Reset machine learning for each domain
	for _, result := range response["results"] {
		rules, ok := result["rule"].([]interface{})
		if !ok {
			resetErrors = append(resetErrors, "FortiWeb returned an invalid rule list")
			continue
		}
		for _, rule := range rules {
			ruleMap, ok := rule.(map[string]interface{})
			if !ok {
				resetErrors = append(resetErrors, "FortiWeb returned an invalid rule")
				continue
			}
			ruleID, hasRuleID := ruleMap["id"]
			domainName, hasDomainName := ruleMap["domain-name"].(string)
			if !hasRuleID || !hasDomainName {
				resetErrors = append(resetErrors, "FortiWeb returned a rule without id or domain-name")
				continue
			}

			resetMLURL := fmt.Sprintf("https://%s:%s/api/v2.0/machine_learning/api_learning_policy.refreshdomain?rule_id=%v", host, port, ruleID)
			// FortiWeb rejects this endpoint when the POST has no HTTP body
			// (errcode -20014). Send an explicit, valid empty JSON object.
			req, err := http.NewRequest("POST", resetMLURL, bytes.NewBufferString("{}"))
			if err != nil {
				errorMsg := fmt.Sprintf("rule_id=%v: unable to create reset request: %v", ruleID, err)
				fmt.Println(errorMsg)
				resetErrors = append(resetErrors, errorMsg)
				continue
			}
			req.Header.Add("Authorization", token)
			req.Header.Add("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				errorMsg := fmt.Sprintf("%s (rule_id=%v): reset request failed: %v", domainName, ruleID, err)
				fmt.Println(errorMsg)
				resetErrors = append(resetErrors, errorMsg)
				continue
			}
			resetBody, readErr := io.ReadAll(resp.Body)
			resp.Body.Close()
			if readErr != nil {
				errorMsg := fmt.Sprintf("%s (rule_id=%v): unable to read reset response: %v", domainName, ruleID, readErr)
				fmt.Println(errorMsg)
				resetErrors = append(resetErrors, errorMsg)
				continue
			}
			if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
				errorMsg := fmt.Sprintf("%s (rule_id=%v): FortiWeb returned %s: %s", domainName, ruleID, resp.Status, strings.TrimSpace(string(resetBody)))
				fmt.Println(errorMsg)
				resetErrors = append(resetErrors, errorMsg)
				continue
			}

			resetDomains = append(resetDomains, domainName)
			fmt.Printf("Machine Learning for domain %s (rule_id=%v) has been reset successfully; FortiWeb status: %s; response: %s\n", domainName, ruleID, resp.Status, strings.TrimSpace(string(resetBody)))
		}
	}

	if len(resetErrors) > 0 {
		return c.String(http.StatusBadGateway, fmt.Sprintf("Machine Learning reset failed: %s", strings.Join(resetErrors, "; ")))
	}
	if len(resetDomains) == 0 {
		return c.String(http.StatusNotFound, "No API Learning rules were found to reset.")
	}

	return c.String(http.StatusOK, fmt.Sprintf("Machine Learning for %s has been reset.", strings.Join(resetDomains, ", ")))
}
