package handlers

import (
	"darwin2/config"
	"darwin2/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os/exec"
	"time"
)

type RequestData struct {
	SelectedOption string `json:"selectedOption"`
}

func HandleWebScan(c echo.Context) error {

	dvwaHost := config.CurrentConfig.DVWAHOST

	// Test if dvwaHost is responding on HTTP port 80
	client := &http.Client{Timeout: 5 * time.Second}

	_, err := client.Get("http://" + dvwaHost)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, fmt.Sprintf("DVWA Host (%s) is not responding on HTTP port 80: %s", dvwaHost, err.Error()))
	}

	_, err = exec.LookPath("nikto")
	if err != nil {
		return c.String(200, "Nikto is not installed on your system")
	}

	var requestData RequestData
	if err := c.Bind(&requestData); err != nil {
		return echo.NewHTTPError(400, "Invalid data")
	}

	randomIP := utils.GenerateRandomPublicIP()

	// Construct the command
	cmd := exec.Command(
		"nikto",
		"-host", dvwaHost,
		"-ask", "no",
		"-followredirects",
		"-maxtime", "60s",
		"-nointeractive",
		"-no404",
		"-timeout", "2",
		"-useragent", "Nikto\r\nX-Forwarded-For: "+randomIP,
		"-T", requestData.SelectedOption)

	// Execute the command and get its output
	output, _ := cmd.CombinedOutput()

	// Return the command output to the client
	return c.String(200, string(output))
}
