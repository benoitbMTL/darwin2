package handlers

import (
	"crypto/tls"
	"darwin2/config"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleKnownBots(c echo.Context) error {
	fmt.Println("Start HandleKnownBot")

	// Retrieve the botName from the form data
	botName := c.FormValue("name")

	// Lookup the userAgent for the given botName
	userAgent, exists := config.BotUserAgents[botName]
	if !exists {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bot '%s' not recognized", botName))
	}

	fmt.Printf("Received botName: %s with userAgent: %s\n", botName, userAgent)

	// Create a new request
	req, err := http.NewRequest("GET", config.CurrentConfig.JUICESHOPURL, nil)
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}

	// Set the User Agent header
	req.Header.Set("User-Agent", userAgent)

	// Create an HTTP client and send the request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	fmt.Println("Reading response body")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return c.String(http.StatusInternalServerError, "Failed to read response body: "+err.Error())
	}

	return c.String(http.StatusOK, string(body))
}
