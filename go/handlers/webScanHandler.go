package handlers

import (
	"crypto/tls"
	"darwin2/config"
	"darwin2/utils"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/labstack/echo/v4"
)

type RequestData struct {
	SelectedOption string `json:"selectedOption"`
}

func HandleWebScan(c echo.Context) error {

	dvwaURL := config.CurrentConfig.DVWAURL

	// Test if DVWA is responding
	client := &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	_, err := client.Get(dvwaURL)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, fmt.Sprintf("DVWA (%s) is not responding: %s", dvwaURL, err.Error()))
	}

	_, err = exec.LookPath("perl")
	if err != nil {
		return c.String(200, "Perl is not installed on your system")
	}

	// Debug: Print the current working directory
	//wd, _ := os.Getwd()
	//fmt.Println("Current working directory:", wd)

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


	var requestData RequestData
	if err := c.Bind(&requestData); err != nil {
		return echo.NewHTTPError(400, "Invalid data")
	}

	randomIP := utils.GenerateRandomPublicIP()

	// Construct the command
	cmd := exec.Command(
		"perl", "nikto/program/nikto.pl",
		"-host", dvwaURL,
		"-ask", "no",
		"-followredirects",
		"-maxtime", "60s",
		"-nointeractive",
		"-no404",
		"-timeout", "2",
		"-useragent", "Nikto Scan Demo\r\nX-Forwarded-For: "+randomIP,
		"-T", requestData.SelectedOption,
	)

	// Execute the command and get its output
	output, _ := cmd.CombinedOutput()

	// Return the command output to the client
	return c.String(200, string(output))
}
