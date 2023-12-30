package handlers

import (
	"crypto/tls"
	"darwin2/config"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// HandleWebAttacks decides which attack to perform based on the request
func HandleWebAttacks(c echo.Context) error {
	attackType := c.FormValue("type")
	username := c.FormValue("username")
	password := config.GetDVWAPassword(username) // Get password for the user

	// Print values to console
	fmt.Printf("Handling web attack - Type: %s, Username: %s, Password: %s\n", attackType, username, password)

	return performAttack(c, attackType, username, password)
}

func performAttack(c echo.Context, attackType, username, password string) error {
	fmt.Println("Performing attack:", attackType)

	// Authenticate and get cookie jar
	cookieJar, err := authenticateUser(username, password)
	if err != nil {
		fmt.Println("Authentication failed:", err)
		return c.String(http.StatusInternalServerError, "Authentication failed: "+err.Error())
	}

	// Get attack configuration
	fmt.Println("Getting attack configuration for:", attackType)
	attackConfig, exists := config.GetAttackConfig(attackType)
	if !exists {
		fmt.Println("Invalid attack type:", attackType)
		return c.String(http.StatusBadRequest, "Invalid attack type")
	}

	// Print attack configuration details
	fmt.Println("Attack Method:", attackConfig.Method)
	fmt.Println("Attack URL:", attackConfig.URL)
	fmt.Println("Attack Post Data:", attackConfig.PostData)

	// Craft the request for the attack
	fmt.Println("Crafting request for attack")
	var req *http.Request
	if attackConfig.Method == "POST" {
		req, err = http.NewRequest(attackConfig.Method, attackConfig.URL, strings.NewReader(attackConfig.PostData))
	} else {
		req, err = http.NewRequest(attackConfig.Method, attackConfig.URL, nil)
	}
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return c.String(http.StatusInternalServerError, "Failed to create request: "+err.Error())
	}

	// Set necessary headers for the attack
	req.Header.Set("User-Agent", config.CurrentConfig.USERAGENT)
	req.Header.Set("Referer", config.CurrentConfig.DVWAURL+"/")
	req.Header.Set("Origin", config.CurrentConfig.DVWAURL)

	// Perform the attack request using the cookie jar
	client := &http.Client{
		Jar: cookieJar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Attack request failed:", err)
		return c.String(http.StatusInternalServerError, "Attack request failed: "+err.Error())
	}
	defer resp.Body.Close()

	// Read the response body
	fmt.Println("Reading response body")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return c.String(http.StatusInternalServerError, "Failed to read response body: "+err.Error())
	}

	// Return the response from the attack
	fmt.Println("Attack performed successfully, response length:", len(body))
	return c.String(http.StatusOK, string(body))
}

// User performs login to the web application and returns the session cookie PHPSESSID
func authenticateUser(username, password string) (*cookiejar.Jar, error) {
	loginURL := config.CurrentConfig.DVWAURL + "/login.php"
	fmt.Println("Login URL:", loginURL)

	formDataString := "username=" + url.QueryEscape(username) + "&password=" + url.QueryEscape(password) + "&Login=Login"
	fmt.Println("Encoded Form Data:", formDataString)

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(formDataString))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", config.CurrentConfig.USERAGENT)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for successful login and return the cookie jar
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed, status: %d", resp.StatusCode)
	}

	return jar, nil
}
