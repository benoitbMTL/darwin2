package handlers

import (
	"crypto/tls"
	"darwin2/config"
	"darwin2/utils"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleBotDeception(c echo.Context) error {
	fmt.Println("Start BotDeception")

	selectedTarget := c.FormValue("selectedTarget")
	var targetURL string
	switch selectedTarget {
	case "DVWA":
		targetURL = config.CurrentConfig.DVWAURL
	case "JuiceShop":
		targetURL = config.CurrentConfig.JUICESHOPURL
	default:
		return c.String(http.StatusBadRequest, "Invalid target selection")
	}

	userAgent := config.CurrentConfig.USERAGENT
	randomIP := utils.GenerateRandomPublicIP()

	// Create a new request
	req, err := http.NewRequest("GET", targetURL+"/fake_url.php", nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the User Agent header
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-FORWARDED-FOR", randomIP)

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
		return fmt.Errorf("error creating request: %v", err)
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

func HandlePageSource(c echo.Context) error {

	selectedTarget := c.FormValue("selectedTarget")
	var targetURL string
	switch selectedTarget {
	case "DVWA":
		targetURL = config.CurrentConfig.DVWAURL
	case "JuiceShop":
		targetURL = config.CurrentConfig.JUICESHOPURL
	default:
		return c.String(http.StatusBadRequest, "Invalid target selection")
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get(targetURL + "/login.php")
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to get page source: %v", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to read page source: %v", err))
	}

	return c.String(http.StatusOK, string(body))
}
