package handlers

import (
	"fmt"
	"os/exec"

	"github.com/labstack/echo/v4"
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

    // Construct the command
    cmd := exec.Command(
        "nikto",
        "-h", "https://192.168.4.40",
        //"-ask", "no",
        //"-Cgidirs", "/cgi/",
        //"-followredirects",
        //"-maxtime", "5s",
        //"-nointeractive",
        //"-no404",
        "-timeout", "2",
        //"-useragent", "Nikto'\n''\r'X-Forwarded-For: 51.13.51.13",
        //"-T", requestData.SelectedOption,
    )

	// Debug print of the command
	fmt.Printf("Executing command: %v\n", cmd.Args)

	// Execute the command and get its output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return echo.NewHTTPError(500, "Error executing nikto command")
	}

	// Debug output of the command
	fmt.Println("Command output:", string(output))

	// Return the command output to the client
	return c.String(200, string(output))
}
