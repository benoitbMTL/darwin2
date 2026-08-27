package handlers

import (
	"crypto/tls"
	"darwin2/config"
	"darwin2/jobs"
	"darwin2/utils"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/labstack/echo/v4"
)

type RequestData struct {
	SelectedTarget string `json:"selectedTarget"`
	SelectedOption string `json:"selectedOption"`
}

func HandleWebScan(c echo.Context) error {
	var requestData RequestData
	if err := c.Bind(&requestData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request data")
	}

	// Debug: Print received data
	fmt.Println("Received scan request for target:", requestData.SelectedTarget, "with option:", requestData.SelectedOption)

	// Map user choice to actual URL
	targetURLMapping := map[string]string{
		"DVWA":      config.CurrentConfig.DVWAURL,
		"Bank":      config.CurrentConfig.BANKURL,
		"JuiceShop": config.CurrentConfig.JUICESHOPURL,
		"Petstore":  config.CurrentConfig.PETSTOREURL,
		"Speedtest": config.CurrentConfig.SPEEDTESTURL,
	}

	targetURL, ok := targetURLMapping[requestData.SelectedTarget]
	if !ok {
		return c.String(http.StatusBadRequest, "Invalid target selection")
	}

	// Test if Target is responding
	client := &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	ctx := c.Request().Context()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, fmt.Sprintf("The Web Server (%s) is not responding: %s", targetURL, err.Error()))
	}
	response.Body.Close()

	_, err = exec.LookPath("perl")
	if err != nil {
		return c.String(200, "Perl is not installed on your system")
	}

	niktoScriptPath := "nikto/program/nikto.pl"

	// Check if the Nikto script exists
	if _, err := os.Stat(niktoScriptPath); err != nil {
		if os.IsNotExist(err) {
			// Return the actual error message if the file does not exist
			return c.String(200, fmt.Sprintf("Nikto is not installed on your system: %s", err.Error()))
		} else {
			// Handle other potential errors from os.Stat
			return c.String(200, fmt.Sprintf("Error checking Nikto installation: %s", err.Error()))
		}
	}

	randomIP := utils.GenerateRandomPublicIP()

	// Construct the command
	jobs.ReportProgress(ctx, 0, 1, "Running Nikto scan")
	cmd := exec.CommandContext(ctx,
		"perl", "nikto/program/nikto.pl",
		"-host", targetURL,
		"-ask", "no",
		"-followredirects",
		"-maxtime", "60s",
		"-nointeractive",
		"-404code", "404",
		"-timeout", "2",
		"-useragent", "Nikto Scan Demo\r\nX-Forwarded-For: "+randomIP,
		"-Tuning", requestData.SelectedOption,
	)

	// Execute the command and get its output
	output, commandErr := cmd.CombinedOutput()
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if commandErr != nil {
		jobs.ReportError(ctx, requestData.SelectedTarget, commandErr)
	}
	jobs.ReportProgress(ctx, 1, 1, "Nikto scan completed")

	// Return the command output to the client
	return c.String(200, string(output))
}
