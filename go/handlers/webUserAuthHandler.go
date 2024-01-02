package handlers

import (
	"crypto/tls"
	"darwin2/config"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func HandleUserAuth(c echo.Context) error {
	username := c.FormValue("username")
	password := config.GetDVWAPassword(username) // Get password for the user
	fmt.Println("Attempting to authenticate user:", username)

	loginURL := config.CurrentConfig.DVWAURL + "/login.php"
	formDataString := "username=" + url.QueryEscape(username) + "&password=" + url.QueryEscape(password) + "&Login=Login"
	fmt.Println("Encoded Form Data for login:", formDataString)

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(formDataString))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return c.String(http.StatusInternalServerError, "Error creating request: "+err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", config.CurrentConfig.USERAGENT)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return c.String(http.StatusInternalServerError, "Error sending request: "+err.Error())
	}
	defer resp.Body.Close()

	fmt.Println("Reading response body from authentication attempt")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return c.String(http.StatusInternalServerError, "Failed to read response body: "+err.Error())
	}

	fmt.Println("Authentication attempt successful, response length:", len(body))
	return c.String(http.StatusOK, string(body))
}
