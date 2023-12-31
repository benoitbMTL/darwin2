package handlers

import (
	"bytes"
	"crypto/tls"
	"darwin2/config"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"strings"
	"net/url"
)

// Struct to map the JSON data returned from api.name-fake.com/random/random
type FakeData struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Birthday  string `json:"birth_data"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	EmailU    string `json:"email_u"`
	EmailD    string `json:"email_d"`
	PhoneH    string `json:"phone_h"`
	Useragent string `json:"useragent"`
	Ipv4      string `json:"ipv4"`
}

    type MLRequest struct {
        SampleCount int `json:"sampleCount"`
    }

func HandleMachineLearning(c echo.Context) error {

    var mlRequest MLRequest

    // Parse request body
    if err := c.Bind(&mlRequest); err != nil {
        return c.String(http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
    }

    // Use the provided sample count or default to 10
    requestCount := mlRequest.SampleCount
    if requestCount <= 0 {
        requestCount = 10
    }

    bankURL := config.CurrentConfig.BANKURL
    var fakeData *FakeData 
    var err error
    
    for i := 0; i < requestCount; i++ {
        log.Printf("Machine Learning: %d\n", i)

        // Fetch random data
        fakeData, err = fetchRandomData()
        if err != nil {
            log.Printf("Error fetching random data: %v\n", err)
            return c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching random data: %v", err))
        }

        // Prepare data for POST request
        postData := preparePostData(*fakeData)

        // Send POST request
        _, err = sendPostRequest(bankURL, postData, fakeData.Useragent, fakeData.Ipv4)
        if err != nil {
            log.Printf("Error sending POST request: %v\n", err)
            return c.String(http.StatusInternalServerError, fmt.Sprintf("Error sending POST request: %v", err))
        }
    }

    // Format the last data for display
    displayData := formatDisplayData(*fakeData)

    // Return a success message with the last data sent
    if requestCount == 1 {
        return c.String(http.StatusOK, fmt.Sprintf("Completed %d request.\n\nLast data sent:\n\n%s", requestCount, displayData))
    } else {
        return c.String(http.StatusOK, fmt.Sprintf("Completed %d requests.\n\nLast data sent:\n\n%s", requestCount, displayData))
    }

}

func fetchRandomData() (*FakeData, error) {

	resp, err := http.Get("https://api.namefake.com/random/random")
	if err != nil {
		log.Printf("Error to Get Random Name: %v\n", err) // Log the error
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body into a byte slice
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data FakeData
	// Use bodyBytes for JSON decoding
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		log.Printf("Error to Decode Random Name: %v\n", err) // Log the error
		return nil, err
	}

	return &data, nil
}

func preparePostData(data FakeData) string {
	names := strings.Split(data.Name, " ")
	firstName := names[0]
	lastName := names[len(names)-1]
	email := fmt.Sprintf("%s@%s", data.EmailU, data.EmailD)
	phone := data.PhoneH

	return fmt.Sprintf("firstname=%s&lastname=%s&email=%s&phone=%s&address=%s&birthday=%s&username=%s&password=%s",
		firstName, lastName, email, phone, url.QueryEscape(data.Address), data.Birthday, data.Username, url.QueryEscape(data.Password))
}

func formatDisplayData(data FakeData) string {
    names := strings.Split(data.Name, " ")
    firstName := names[0]
    lastName := names[len(names)-1]
    email := fmt.Sprintf("%s@%s", data.EmailU, data.EmailD)
    phone := data.PhoneH

    // Sanitize the address
    sanitizedAddress := strings.ReplaceAll(data.Address, "\n", " ") // Replace newline characters
    sanitizedAddress = strings.ReplaceAll(sanitizedAddress, "\t", " ") // Replace tab characters

    return fmt.Sprintf(
        "Firstname:\t%s\nLastname:\t%s\nEmail:\t\t%s\nPhone:\t\t%s\nAddress:\t%s\nBirthday:\t%s\nUsername:\t%s\nPassword:\t%s",
        firstName, lastName, email, phone, sanitizedAddress, data.Birthday, data.Username, data.Password)
}


func sendPostRequest(url, data, userAgent, ipv4 string) (*http.Response, error) {

	// Create a custom http.Transport with TLSClientConfig
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Use the custom Transport with http.Client
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "http://bank.corp.fabriclab.ca")
	req.Header.Set("Referer", "http://bank.corp.fabriclab.ca/bank.html")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-Forwarded-For", ipv4)

	return client.Do(req)
}
