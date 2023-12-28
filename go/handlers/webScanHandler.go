package handlers

import (
	//"fmt"
	"darwin2/utils"
	"github.com/labstack/echo/v4"
	"os/exec"
)

type RequestData struct {
	SelectedOption string `json:"selectedOption"`
}

func HandleWebScan(c echo.Context) error {
	_, err := exec.LookPath("nikto")
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
		"-host", "192.168.4.40",
		"-ask", "no",
		"-followredirects",
		"-maxtime", "60s",
		"-nointeractive",
		"-no404",
		"-timeout", "2",
		"-useragent", "Nikto\r\nX-Forwarded-For: "+randomIP,
		"-T", requestData.SelectedOption)

	//cmd.Dir = "/tmp" // Set the working directory to /tmp

	// Debug print of the command
	//fmt.Printf("Executing command: %v\n", cmd)

	// Execute the command and get its output
	output, _ := cmd.CombinedOutput()

	// Return the command output to the client
	return c.String(200, string(output))
}
