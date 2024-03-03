package handlers

import (
	"crypto/tls"
	"darwin2/config"
	"darwin2/utils"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleBotDeception(c echo.Context) error {
	fmt.Println("Start BotDeception")

	userAgent := config.CurrentConfig.USERAGENT
	randomIP := utils.GenerateRandomPublicIP()

	// Create a new request
	req, err := http.NewRequest("GET", config.CurrentConfig.DVWAURL+"/fake_url.php", nil)
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
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


func HandlePageSource(c echo.Context) error {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get(config.CurrentConfig.DVWAURL + "/login.php")
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to get page source: %v", err))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to read page source: %v", err))
	}

	return c.String(http.StatusOK, string(body))
}


