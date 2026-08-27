package handlers

import (
	"context"
	"crypto/tls"
	"darwin2/config"
	"darwin2/jobs"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
)

type HealthCheckResult struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	URL                  string `json:"url"`
	Status               string `json:"status"`
	HTTPCode             int    `json:"httpCode,omitempty"`
	DurationMilliseconds int64  `json:"durationMilliseconds"`
	Message              string `json:"message,omitempty"`
}

type HealthCheckSummary struct {
	Up            int `json:"up"`
	Down          int `json:"down"`
	NotConfigured int `json:"notConfigured"`
	Total         int `json:"total"`
}

type HealthCheckResponse struct {
	Status               string              `json:"status"`
	CheckedAt            time.Time           `json:"checkedAt"`
	DurationMilliseconds int64               `json:"durationMilliseconds"`
	Summary              HealthCheckSummary  `json:"summary"`
	Results              []HealthCheckResult `json:"results"`
}

type healthTarget struct {
	id      string
	name    string
	url     string
	apiTest bool
	token   string
}

func HandleHealthCheck(c echo.Context) error {
	startedAt := time.Now()
	currentConfig := config.GetCurrentConfig()
	tokenData := fmt.Sprintf(`{"username":"%s","password":"%s","vdom":"%s"}`, currentConfig.USERNAMEAPI, currentConfig.PASSWORDAPI, currentConfig.VDOMAPI)
	apiToken := base64.StdEncoding.EncodeToString([]byte(tokenData))
	managementURL := ""
	if currentConfig.FWBMGTIP != "" && currentConfig.FWBMGTPORT != "" {
		managementURL = fmt.Sprintf("https://%s:%s", currentConfig.FWBMGTIP, currentConfig.FWBMGTPORT)
	}

	targets := []healthTarget{
		{"dvwa", "DVWA", currentConfig.DVWAURL, false, ""},
		{"bank", "Bank", currentConfig.BANKURL, false, ""},
		{"juice-shop", "Juice Shop", currentConfig.JUICESHOPURL, false, ""},
		{"petstore", "Petstore API", currentConfig.PETSTOREURL, false, ""},
		{"speedtest", "Speedtest", currentConfig.SPEEDTESTURL, false, ""},
		{"fortiguard", "FortiGuard", "https://www.fortiguard.com", false, ""},
		{"fortiweb-management", "FortiWeb Management", managementURL, false, ""},
		{"fortiweb-api", "FortiWeb API", managementURL + "/api/v2.0/cmdb/system/global", true, apiToken},
	}

	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	results := make([]HealthCheckResult, len(targets))
	var waitGroup sync.WaitGroup
	var completed atomic.Int64
	for index, target := range targets {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			results[index] = checkHealthTarget(c.Request().Context(), client, target)
			current := completed.Add(1)
			jobs.ReportProgress(c.Request().Context(), current, int64(len(targets)), fmt.Sprintf("Checked %d of %d resources", current, len(targets)))
		}()
	}
	waitGroup.Wait()

	summary := HealthCheckSummary{Total: len(results)}
	for _, result := range results {
		switch result.Status {
		case "up":
			summary.Up++
		case "not-configured":
			summary.NotConfigured++
		default:
			summary.Down++
		}
	}
	status := "healthy"
	if summary.Down > 0 {
		status = "degraded"
	}

	return c.JSON(http.StatusOK, HealthCheckResponse{
		Status:               status,
		CheckedAt:            time.Now().UTC(),
		DurationMilliseconds: time.Since(startedAt).Milliseconds(),
		Summary:              summary,
		Results:              results,
	})
}

func checkHealthTarget(ctx context.Context, client *http.Client, target healthTarget) HealthCheckResult {
	result := HealthCheckResult{ID: target.id, Name: target.name, URL: target.url}
	if target.url == "" || (target.apiTest && target.url == "/api/v2.0/cmdb/system/global") {
		result.Status = "not-configured"
		result.Message = "Not configured"
		return result
	}

	startedAt := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target.url, nil)
	if err != nil {
		result.Status = "down"
		result.Message = err.Error()
		return result
	}
	if target.apiTest {
		req.Header.Set("Authorization", target.token)
		req.Header.Set("Accept", "application/json")
	}

	response, err := client.Do(req)
	result.DurationMilliseconds = time.Since(startedAt).Milliseconds()
	if err != nil {
		result.Status = "down"
		result.Message = err.Error()
		return result
	}
	defer response.Body.Close()
	result.HTTPCode = response.StatusCode

	if target.apiTest {
		var payload struct {
			Results struct {
				Hostname string `json:"hostname"`
			} `json:"results"`
		}
		if err := json.NewDecoder(response.Body).Decode(&payload); err != nil || payload.Results.Hostname == "" {
			result.Status = "down"
			result.Message = "API configuration is incorrect"
			return result
		}
		result.Status = "up"
		result.Message = payload.Results.Hostname
		return result
	}

	result.Status = "up"
	result.Message = response.Status
	return result
}
