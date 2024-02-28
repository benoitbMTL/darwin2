package handlers

import (
	"crypto/tls"
	"darwin2/config"
	"darwin2/utils"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type RequestData struct {
	SelectedOption string `json:"selectedOption"`
}

func HandleWebScan(c echo.Context) error {

	dvwaHost := config.CurrentConfig.DVWAHOST

	// Test if dvwaHost is responding on HTTP port 443
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Test if dvwaHost is responding on HTTPS port 443
	_, err := client.Get("https://" + dvwaHost)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, fmt.Sprintf("DVWA Host (%s) is not responding on HTTPS port 443: %s", dvwaHost, err.Error()))
	}

	_, err = exec.LookPath("perl")
	if err != nil {
		return c.String(200, "Perl is not installed on your system")
	}

	niktoScriptPath := "nikto/program/nikto.pl"

	// Check if the Nikto script exists
	if _, err := os.Stat(niktoScriptPath); os.IsNotExist(err) {
		return c.String(200, "Nikto is not installed on your system")
	}

	var requestData RequestData
	if err := c.Bind(&requestData); err != nil {
		return echo.NewHTTPError(400, "Invalid data")
	}

	randomIP := utils.GenerateRandomPublicIP()

	// Construct the command
	cmd := exec.Command(
		"perl nikto/program/nikto.pl",
		"-host", dvwaHost,
		"-ask", "no",
		"-followredirects",
		"-maxtime", "60s",
		"-nointeractive",
		"-no404",
		"-timeout", "2",
		"-useragent", "Nikto Scan Demo\r\nX-Forwarded-For: "+randomIP,
		"-T", requestData.SelectedOption)

	// Execute the command and get its output
	output, _ := cmd.CombinedOutput()

	// Return the command output to the client
	return c.String(200, string(output))
}
